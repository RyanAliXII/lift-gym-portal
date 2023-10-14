import { createApp, ref } from "vue";

createApp({
  setup() {
    const form = ref({
      id: 0,
      description: "",
      price: 0,
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
    const errors = ref({});

    return {
      form,
      errors,
      handleFormInput,
    };
  },
}).mount("#CoachingRate");
