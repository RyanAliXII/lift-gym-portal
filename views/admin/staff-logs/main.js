import Choices from "choices.js";
import { createApp, onMounted, ref } from "vue";
import { useDebounceFn } from "@vueuse/core";
import DataTable from "datatables.net-vue3";
import DataTableCore from "datatables.net";
import swal from "sweetalert2";
import "datatables.net-dt/css/jquery.dataTables.min.css";

DataTable.use(DataTableCore);
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  components: {
    DataTable,
  },
  setup() {
    const initialForm = {
      id: 0,
      staffId: 0,
    };

    const table = ref(null);
    let dt;
    const columns = [
      {
        title: "In",
        data: "createdAt",
        render: (value) => {
          return `<span>${formatDate(value)}</span>`;
        },
      },
      {
        title: "Out",
        data: "isLoggedOut",
        render: (value, _, row) => {
          if (value) {
            return `<span>${formatDate(row.loggedOutAt)}</span>`;
          }

          return "Not logged out.";
        },
      },
      {
        title: "Staff ID",
        data: "staff.publicId",
      },
      {
        title: "Staff",
        data: null,
        render: (value, event, row) =>
          `${row.staff.givenName} ${row.staff.surname}`,
      },
      {
        title: "",
        data: null,
        render: (value, event, row) => {
          let buttons = `<div class='d-flex' style='gap:5px'>`;
          if (window.hasEditPermission) {
            buttons += `<button
            class="btn btn-outline-primary edit-log"
            data-toggle="tooltip"
            title="Edit Log"
          
          >
            <i class="fas fa-edit"></i>
          </button>`;
          }

          if (window.hasEditPermission && !row.isLoggedOut) {
            buttons += `<button
            class="btn btn-outline-secondary logout-btn"
            data-toggle="tooltip"
            title="Logout"
          
          >
          <i class="fas fa-sign-out-alt"></i>
          </button>`;
          }
          if (window.hasDeletePermission) {
            buttons += `
            <button
              data-id=${row.id}
              class="btn btn-outline-danger delete-log"
              data-toggle="tooltip"
              title="Delete Log"
         
          >
            <i class="fas fa-trash"></i>
          </button>
            `;
          }

          return buttons + `</div>`;
        },
      },
    ];
    const tableConfig = {
      lengthMenu: [25],
      lengthChange: false,
      dom: "lrtip",
    };
    const logs = ref([]);
    const logClientSelectElement = ref(null);
    const logClientSelect = ref(null);
    const editLogClientSelectElement = ref(null);
    const editLogClientSelect = ref(null);
    const form = ref({
      ...initialForm,
    });
    const errors = ref({});
    const handleFormInput = (event) => {
      let value = event.target.value;
      let name = event.target.name;
      if (event.target.type === "number") {
        value = Number(value);
      }
      form.value[name] = value;
      delete errors.value[name];
    };
    const fetchLogs = async () => {
      try {
        const response = await fetch("/app/staff-logs", {
          headers: new Headers({ "Content-Type": "application/json" }),
        });
        if (response.status >= 400) return;
        const { data } = await response.json();
        logs.value = data?.staffLogs ?? [];
      } catch (error) {
        console.error(error);
      }
    };
    const fetchStaffByKeyword = async (query) => {
      const response = await fetch(
        `/app/staffs?${new URLSearchParams({
          keyword: query,
        }).toString()}`,
        {
          headers: new Headers({
            "Content-Type": "application/json",
            "Cache-Control": "no-cache",
          }),
        }
      );

      if (response.status === 200) {
        const { data } = await response.json();
        const selectValues = (data?.staffs ?? []).map((staff) => ({
          value: staff.id,
          label: `${staff.givenName} ${staff.surname} - ${staff.email}`,
          customProperties: staff,
        }));
        logClientSelect.value.setChoices(selectValues, "value", "label", true);
        editLogClientSelect.value.setChoices(
          selectValues,
          "value",
          "label",
          true
        );
      }
    };
    const search = useDebounceFn(fetchStaffByKeyword, 500);
    const submitLog = async () => {
      errors.value = {};
      try {
        const response = await fetch("/app/staff-logs", {
          method: "POST",
          body: JSON.stringify(form.value),
          headers: new Headers({
            "Content-Type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
        });
        const { data } = await response.json();
        if (response.status >= 400) {
          if (data?.errors) {
            errors.value = data?.errors;
          }
          return;
        }
        swal.fire(
          "Staff Loggged In",
          "Staff has been loggged in successfully",
          "success"
        );
        form.value = {
          ...initialForm,
        };
        fetchLogs();
        $("#logClientModal").modal("hide");
      } catch (error) {
        console.error(error);
      }
    };
    const updateLog = async () => {
      errors.value = {};
      try {
        const response = await fetch(`/app/staff-logs/${form.value.id}`, {
          method: "PUT",
          body: JSON.stringify(form.value),
          headers: new Headers({
            "Content-Type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
        });
        const { data } = await response.json();
        if (response.status >= 400) {
          if (data?.errors) {
            errors.value = data?.errors;
          }
          return;
        }
        swal.fire(
          "Staff Log Updated",
          "Staff log has been updated.",
          "success"
        );
        form.value = {
          ...initialForm,
        };
        $("#editLogModal").modal("hide");
        fetchLogs();
      } catch (error) {
        console.error(error);
      }
    };
    const deleteLog = async (id) => {
      errors.value = {};
      try {
        const response = await fetch(`/app/staff-logs/${id}`, {
          method: "DELETE",
          headers: new Headers({
            "Content-Type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
        });
        const { data } = await response.json();
        if (response.status >= 400) {
          if (data?.errors) {
            errors.value = data?.errors;
          }
          return;
        }
        swal.fire(
          "Staff Log Deleted",
          "Staff log has been deleted.",
          "success"
        );
        fetchLogs();
      } catch (error) {
        console.error(error);
      }
    };
    const logoutClient = async (id) => {
      errors.value = {};
      try {
        const response = await fetch(`/app/staff-logs/${id}/logout`, {
          method: "PATCH",
          headers: new Headers({
            "Content-Type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
        });
        const { data } = await response.json();
        if (response.status >= 400) {
          if (data?.errors) {
            errors.value = data?.errors;
          }
          return;
        }
        swal.fire(
          "Staff Logged Out",
          "Staff log has been logged out.",
          "success"
        );
        fetchLogs();
      } catch (error) {
        console.error(error);
      }
    };
    const initModalListeners = () => {
      $("#logClientModal").on("hidden.bs.modal", () => {
        logClientSelect.value.removeActiveItems();
        form.value = {
          ...initialForm,
        };
      });
      $("#editLogModal").on("hidden.bs.modal", () => {
        editLogClientSelect.value.removeActiveItems();
        form.value = {
          ...initialForm,
        };
      });
    };

    const formatDate = (date) => {
      return new Date(date).toLocaleDateString("en-US", {
        month: "long",
        day: "2-digit",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
      });
    };
    const toMoney = (money) => {
      if (!money) return 0;
      return money.toLocaleString(undefined, {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
      });
    };
    onMounted(() => {
      dt = table?.value?.dt;
      logClientSelect.value = new Choices(logClientSelectElement.value, {
        allowHTML: false,
        placeholder: "Seach Client",
      });
      editLogClientSelect.value = new Choices(
        editLogClientSelectElement.value,
        {
          allowHTML: false,
          placeholder: "Seach Client",
        }
      );

      logClientSelect.value.passedElement.element.addEventListener(
        "search",
        (event) => {
          search(event.detail.value);
        }
      );
      editLogClientSelect.value.passedElement.element.addEventListener(
        "search",
        (event) => {
          search(event.detail.value);
        }
      );
      logClientSelect.value.passedElement.element.addEventListener(
        "change",
        () => {
          const select = logClientSelect.value.getValue();
          form.value = {
            ...form.value,
            staffId: select.value,
          };
          delete errors.coachId;
        }
      );

      editLogClientSelect.value.passedElement.element.addEventListener(
        "change",
        () => {
          const select = editLogClientSelect.value.getValue();

          if (!select) return;

          form.value = {
            ...form.value,
            staffId: select.value,
          };
          delete errors.clientId;
        }
      );
      initModalListeners();
      fetchLogs();

      $(dt.table().body()).on("click", "button.delete-log", async (event) => {
        const id = event.currentTarget.getAttribute("data-id");
        initDelete(id);
      });
      $(dt.table().body()).on("click", "button.edit-log", async (event) => {
        let data = dt.row(event.target.closest("tr")).data();
        initEdit(data);
      });

      $(dt.table().body()).on("click", "button.logout-btn", async (event) => {
        let data = dt.row(event.target.closest("tr")).data();
        const result = await swal.fire({
          showCancelButton: true,
          confirmButtonText: "Yes, logout staff.",
          title: "Staff Logout",
          text: "Are you sure you want to logout staff?",
          confirmButtonColor: "#295ad6",
          cancelButtonText: "I don't want to logout staff.",
          icon: "question",
        });
        if (result.isConfirmed) {
          logoutClient(data.id);
        }
      });
    });

    const searchLogs = (event) => {
      const query = event.target.value;
      dt.search(query).draw();
    };
    const initEdit = async (log) => {
      form.value = {
        id: log.id,
        staffId: log.staffId,
      };
      await fetchStaffByKeyword(log.staff.publicId);
      editLogClientSelect.value.setChoiceByValue(log.staff.id);
      $("#editLogModal").modal("show");
    };
    const initDelete = async (id) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, delete it",
        title: "Delete Staff Log",
        text: "Are you sure you want to delete staff log?",
        confirmButtonColor: "#d9534f",
        cancelButtonText: "I don't want to delete the log",
        icon: "warning",
      });
      if (result.isConfirmed) {
        deleteLog(id);
      }
    };

    return {
      logClientSelectElement,
      form,
      logs,
      formatDate,
      handleFormInput,
      submitLog,
      errors,
      toMoney,
      initEdit,
      updateLog,
      editLogClientSelectElement,
      initDelete,
      columns,
      table,
      tableConfig,
      searchLogs,
    };
  },
}).mount("#ClientLog");
