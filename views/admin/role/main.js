const { createApp, onMounted, ref } = require("vue");
import { useForm } from "@vorms/core";
import Choices from "choices.js";
createApp({
  setup() {
    const addSelectElement = ref(null);
    const addSelect = ref(null);
    const { register, handleSubmit, handleReset, errors } = useForm({
      initialValues: {
        name: "",
        permissions: [],
      },
      onSubmit(data) {
        console.log(data);
      },
    });

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
    };
  },
}).mount("#RolePage");
