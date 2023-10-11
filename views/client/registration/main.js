import { useForm } from "vee-validate";
import { createApp, ref } from "vue";
import { object } from "yup";

createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const displaySuccessMessage = ref(false);
    const {
      defineInputBinds,
      handleSubmit,
      setErrors,
      errors,
      isSubmitting,
      resetForm,
    } = useForm({
      initialValues: {
        givenName: "",
        middleName: "",
        surname: "",
        email: "",
        password: "",
        dateOfBirth: "",
      },
      validationSchema: object({}),
      validateOnMount: false,
    });

    const givenName = defineInputBinds("givenName", { validateOnInput: true });
    const middleName = defineInputBinds("middleName", {
      validateOnInput: true,
    });
    const surname = defineInputBinds("surname", { validateOnInput: true });
    const email = defineInputBinds("email", { validateOnInput: true });
    const password = defineInputBinds("password", { validateOnInput: true });
    const dateOfBirth = defineInputBinds("dateOfBirth", {
      validateOnInput: true,
    });
    const onSubmit = handleSubmit(async (values) => {
      const response = await fetch("/clients/registration", {
        method: "POST",
        body: JSON.stringify(values),
        headers: new Headers({
          "Content-Type": "application/json",
          "X-CSRF-Token": window.csrf,
        }),
      });
      const { data } = await response.json();

      if (response.status >= 400) {
        if (data?.errors) {
          setErrors({ ...data?.errors });
        }
        return;
      }
      displaySuccessMessage.value = true;
      resetForm();
    });

    return {
      givenName,
      surname,
      email,
      password,
      dateOfBirth,
      middleName,
      onSubmit,
      isSubmitting,
      errors,
      displaySuccessMessage,
    };
  },
}).mount("#RegistrationPage");
