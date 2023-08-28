import { useForm } from "vee-validate";
import { createApp, onMounted } from "vue";

createApp({
  setup() {
    const { errors, defineInputBinds, values } = useForm({
      initialValues: {
        givenName: "",
        middleName: "",
        surname: "",
        dateOfBirth: "",
        address: "",
        email: "",
        mobileNumber: "",
        emergencyContact: "",
      },
      validateOnMount: false,
    });

    const givenName = defineInputBinds("givenName");
    const middleName = defineInputBinds("middleName");
    const surname = defineInputBinds("surname");
    const dateOfBirth = defineInputBinds("dateOfBirth");
    const address = defineInputBinds("address");
    const email = defineInputBinds("email");
    const mobileNumber = defineInputBinds("mobileNumber");
    const emergencyContact = defineInputBinds("emergencyContact");

    const onSubmit = () => {
      console.log(values);
    };
    return {
      givenName,
      middleName,
      surname,
      dateOfBirth,
      address,
      email,
      mobileNumber,
      emergencyContact,
      onSubmit,
      errors,
    };
  },
}).mount("#CoachRegistrationPage");
