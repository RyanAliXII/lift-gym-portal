import { createApp, ref } from "vue";

createApp({
  setup() {
    const form = ref({ email: "" });
    const handleFormInput = (event) => {
      let value = event.target.value;
      let name = event.target.name;
      if (event.target.type === "number") {
        value = Number(value);
      }
      form.value[name] = value;
    };
    return {
      form,
      handleFormInput,
    };
  },
}).mount("#ResetPassword");
