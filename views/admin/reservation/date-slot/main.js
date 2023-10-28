import { format } from "date-fns";
import { createApp, ref } from "vue";

createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const today = format(new Date(), "yyyy-MM-dd");
    const initialFormValue = {
      to: today,
      from: today,
    };
    const form = ref({
      ...initialFormValue,
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
    const onSubmit = async () => {
      try {
        errors.value = {};
        const response = await fetch("/app/date-slots", {
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
            errors.value = data.errors;
          }
          return;
        }
        form.value = { ...initialFormValue };
        $("#newSlotModal").modal("hide");
      } catch (error) {
        console.error(error);
      }
    };
    return {
      form,
      handleFormInput,
      onSubmit,
      errors,
    };
  },
}).mount("#DateSlot");
