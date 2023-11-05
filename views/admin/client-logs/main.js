import Choices from "choices.js";
import { createApp, onMounted, ref } from "vue";
import { useDebounce, useDebounceFn } from "@vueuse/core";
import swal from "sweetalert2";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const initialForm = {
      id: 0,
      clientId: 0,
      isMember: false,
      amountPaid: 0,
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
        const response = await fetch("/app/client-logs", {
          headers: new Headers({ "Content-Type": "application/json" }),
        });
        if (response.status >= 400) return;
        const { data } = await response.json();
        logs.value = data?.clientLogs ?? [];
      } catch (error) {
        console.error(error);
      }
    };
    const fetchClientByKeyword = async (query) => {
      const response = await fetch(
        `/app/clients?${new URLSearchParams({
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
        const selectValues = (data?.clients ?? []).map((client) => ({
          value: client.id,
          label: `${client.givenName} ${client.surname} - ${client.email} - ${
            client.isMember ? "Member" : "Non-Member"
          }`,
          customProperties: client,
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
    const search = useDebounceFn(fetchClientByKeyword, 500);
    const submitLog = async () => {
      errors.value = {};
      try {
        const response = await fetch("/app/client-logs", {
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
          "Client Loggged In",
          "Client has been loggged in successfully",
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
        const response = await fetch(`/app/client-logs/${form.value.id}`, {
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
        swal.fire("Client Log Updated", "Client has been updated.", "success");
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
        const response = await fetch(`/app/client-logs/${id}`, {
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
          "Client Log Deleted",
          "Client log has been deleted.",
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
            clientId: select.value,
            isMember: select.customProperties.isMember,
          };
          delete errors.clientId;
        }
      );

      editLogClientSelect.value.passedElement.element.addEventListener(
        "change",
        () => {
          const select = editLogClientSelect.value.getValue();

          if (!select) return;

          form.value = {
            ...form.value,
            clientId: select.value,
            isMember: select.customProperties.isMember,
          };
          delete errors.clientId;
        }
      );
      initModalListeners();
      fetchLogs();
    });
    const initEdit = async (log) => {
      form.value = {
        id: log.id,
        clientId: log.clientId,
        amountPaid: log.amountPaid,
        isMember: log.isMember,
      };
      await fetchClientByKeyword(log.client.givenName);
      editLogClientSelect.value.setChoiceByValue(log.clientId);
      $("#editLogModal").modal("show");
    };
    const initDelete = async (id) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, delete it",
        title: "Delete Client Log",
        text: "Are you sure you want to delete client log?",
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
    };
  },
}).mount("#ClientLog");
