import { useForm } from "vee-validate";
import { createApp, ref } from "vue";
import { object } from "yup";

createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const displaySuccessMessage = ref(false);
    const showPassword = ref(false);
    const isSubmitting = ref(false);
    const form = ref({
      givenName: "",
      middleName: "",
      surname: "",
      gender: "",
      email: "",
      password: "",
      dateOfBirth: "",
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
      isSubmitting.value = true;
      const response = await fetch("/clients/registration", {
        method: "POST",
        body: JSON.stringify(form.value),
        headers: new Headers({
          "Content-Type": "application/json",
          "X-CSRF-Token": window.csrf,
        }),
      });
      isSubmitting.value = false;
      const { data } = await response.json();

      if (response.status >= 400) {
        if (data?.errors) {
          errors.value = data?.errors;
        }

        return;
      }
      displaySuccessMessage.value = true;
      resetForm();
    };

    return {
      errors,
      onSubmit,
      form,
      displaySuccessMessage,
      isSubmitting,
      handleFormInput,
      showPassword,
    };
  },
}).mount("#RegistrationPage");
