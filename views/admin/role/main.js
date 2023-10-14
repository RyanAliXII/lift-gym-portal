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
      const name = event.target.name;
      errors.value = { ...errors.value, name: undefined };
      form.value = { ...form, [name]: event.target.value };
    };

    const fetchRoles = async () => {
      try {
        const response = await fetch("/app/roles", {
          headers: new Headers({ "Content-Type": "application/json" }),
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
        swal.fire("New Role", "Role has been created.", "success");
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
    return {
      addSelectElement,
      form,
      handleFormInput,
      errors,
      roles,
      initEdit,
      editSelectElement,
      onSubmitNew,
      onSubmitUpdate,
    };
  },
}).mount("#RolePage");
