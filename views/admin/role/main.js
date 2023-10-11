const { createApp, onMounted, ref } = require("vue");

import Choices from "choices.js";
import swal from "sweetalert2";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const addSelectElement = ref(null);
    const addSelect = ref(null);
    const roles = ref([]);
    const form = ref({
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
    });
    return {
      addSelectElement,
      form,
      handleFormInput,
      errors,
      roles,
      onSubmitNew,
    };
  },
}).mount("#RolePage");
