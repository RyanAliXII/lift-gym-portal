const { createApp, onMounted, ref } = require("vue");
import { useForm } from "@vorms/core";
import Choices from "choices.js";
import { add } from "date-fns";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const addSelectElement = ref(null);
    const addSelect = ref(null);
    const { register, handleSubmit, handleReset, values, errors, setValues } =
      useForm({
        initialValues: {
          name: "",
          permissions: [],
        },
        async onSubmit(data, event) {
          const permissions = addSelect.value.getValue().map((p) => p.value);
          await submitNewRole({ ...data, permissions: permissions });
        },
      });
    const submitNewRole = async (form) => {
      try {
        const response = await fetch("/app/roles", {
          method: "POST",
          body: JSON.stringify(form),
          headers: new Headers({
            "Content-Type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
        });
      } catch (err) {
        console.error(err);
      }
    };
    const { value: name, attrs: nameAttrs } = register("name");

    const { value: permissions, attrs: permissionAttrs } =
      register("permissions");

    onMounted(() => {
      addSelect.value = new Choices(addSelectElement.value, {
        allowHTML: true,
      });
    });
    return {
      addSelectElement,
      name,
      permissions,
      errors,
      handleSubmit,
    };
  },
}).mount("#RolePage");
