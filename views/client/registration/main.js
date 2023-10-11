import { format } from "date-fns";
import { useForm } from "vee-validate";
import { createApp, onMounted } from "vue";
import { object } from "yup";

createApp({
  setup() {
    const { defineInputBinds, handleSubmit } = useForm({
      initialValues: {
        givenName: "",
        surname: "",
        email: "",
        password: "",
        dateOfBirth: format(new Date(), "yyyy-MM-dd"),
      },

      validateOnMount: false,
      validationSchema: object({}),
    });

    const givenName = defineInputBinds("givenName", { validateOnChange: true });
    const surname = defineInputBinds("surname", { validateOnChange: true });
    const email = defineInputBinds("email", { validateOnChange: true });
    const password = defineInputBinds("password", { validateOnChange: true });
    const repeatPasword = defineInputBinds("repeatPasword", {
      validateOnChange: true,
    });
    const dateOfBirth = defineInputBinds("dateOfBirth", {
      validateOnChange: true,
    });
    const onSubmit = handleSubmit(async (values) => {});
    return {
      givenName,
      surname,
      email,
      password,
      dateOfBirth,
      onSubmit,
    };
  },
}).mount("#RegistrationPage");
