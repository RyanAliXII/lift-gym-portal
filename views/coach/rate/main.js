import { createApp, ref } from "vue";

createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
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

    const onSubmitNewRate = async () => {
      try {
        const response = await fetch("/coaches/rates", {
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
            errors.value = data?.errors ?? {};
          }
          return;
        }
      } catch (error) {
        console.error(error);
      }
    };

    return {
      form,
      errors,
      handleFormInput,
      onSubmitNewRate,
    };
  },
}).mount("#CoachingRate");
