const { createApp, onMounted, ref } = require("vue");

import Choices from "choices.js";
import swal from "sweetalert2";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const addSelectElement = ref(null);
    const editSelectElement = ref(null);
    const addSelect = ref(null);
    const editSelect = ref(null);
    const roles = ref([]);
    const form = ref({
      id: 0,
      name: "",
      permissions: [],
    });
    const errors = ref({
      name: "",
      permissions: "",
    });

    const handleFormInput = (event) => {
      let value = event.target.value;
      let name = event.target.name;
      if (event.target.type === "number") {
        value = Number(value);
      }
      form.value[name] = value;
      delete errors.value[name];
    };

    const fetchRoles = async () => {
      try {
        const response = await fetch("/app/roles", {
          headers: new Headers({
            "Content-Type": "application/json",
            "Cache-Control": "no-cache",
          }),
        });
        const { data } = await response.json();
        roles.value = data?.roles ?? [];
      } catch (err) {
        roles.value = [];
      }
    };
    const onSubmitUpdate = async () => {
      try {
        errors.value = {};
        const permissions = editSelect.value.getValue().map((p) => p.value);
        const response = await fetch(`/app/roles/${form.value.id}`, {
          method: "PUT",
          body: JSON.stringify({ ...form.value, permissions: permissions }),
          headers: new Headers({
            "Content-Type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
        });
        const { data } = await response.json();
        if (response.status >= 400) {
          if (data?.errors) {
            errors.value = data.errors;
          }
          return;
        }
        form.value = {
          id: 0,
          name: "",
          permissions: [],
        };
        fetchRoles();
        editSelect.value.removeActiveItems();
        $("#editRoleModal").modal("hide");
        swal.fire("New Role", "Role has updated.", "success");
      } catch (err) {
        console.error(err);
      }
    };

    const onSubmitNew = async () => {
      try {
        errors.value = {};
        const permissions = addSelect.value.getValue().map((p) => p.value);
        const response = await fetch("/app/roles", {
          method: "POST",
          body: JSON.stringify({ ...form.value, permissions: permissions }),
          headers: new Headers({
            "Content-Type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
        });
        const { data } = await response.json();
        if (response.status >= 400) {
          if (data?.errors) {
            errors.value = data.errors;
          }
          return;
        }
        form.value = {
          name: "",
          permissions: [],
        };
        fetchRoles();
        addSelect.value.removeActiveItems();
        $("#newRoleModal").modal("hide");
        swal.fire("New Role", "Role has been created.", "success");
      } catch (err) {
        console.error(err);
      }
    };

    onMounted(() => {
      fetchRoles();
      addSelect.value = new Choices(addSelectElement.value, {
        allowHTML: true,
      });
      editSelect.value = new Choices(editSelectElement.value, {
        allowHTML: true,
      });
      $("#newRoleModal").on("hidden.bs.modal", () => {
        form.value = {
          name: "",
          permissions: [],
        };
        addSelect.value.removeActiveItems();
      });
      $("#editRoleModal").on("hidden.bs.modal", () => {
        form.value = {
          name: "",
          permissions: [],
        };
        editSelect.value.removeActiveItems();
      });
    });

    const initEdit = (role) => {
      form.value = { ...role };
      editSelect.value.setChoiceByValue(role.permissions);
      $("#editRoleModal").modal("show");
    };
    const initDelete = async (id) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, delete it",
        title: "Delete Role",
        text: "Are you sure you want to delete role?",
        confirmButtonColor: "#d9534f",
        cancelButtonText: "I don't want to delete this role",
        icon: "warning",
      });
      if (result.isConfirmed) {
        deleteRole(id);
      }
    };
    const deleteRole = async (id) => {
      const response = await fetch(`/app/roles/${id}`, {
        method: "DELETE",
        headers: new Headers({
          "Content-Type": "application/json",
          "X-CSRF-Token": window.csrf,
        }),
      });

      if (response.status === 200) {
        swal.fire("Delete Role", "Role has been deleted.", "success");
        fetchRoles();
      }
    };
    return {
      addSelectElement,
      form,
      handleFormInput,
      errors,
      roles,
      initEdit,
      initDelete,
      editSelectElement,
      onSubmitNew,
      onSubmitUpdate,
    };
  },
}).mount("#RolePage");
