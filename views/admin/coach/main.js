import { useForm, validate } from "vee-validate";
import { createApp } from "vue";
import swal from "sweetalert2";
import { object } from "yup";
createApp({
  setup() {
    const { errors, defineInputBinds, values, setErrors } = useForm({
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
      validationSchema: object({}),
      validateOnMount: false,
    });
    const newCoach = async () => {
      try {
        const response = await fetch("/app/coaches", {
          headers: new Headers({
            "content-type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
          method: "POST",
          body: JSON.stringify(values),
        });

        const { data } = await response.json();
        if (response.status === 200) {
          swal.fire(
            "Coach Register",
            "Coach has been registered successfully.",
            "success"
          );
        }
        if (response.status === 400 && data?.errors) {
          setErrors(data.errors);
        }
      } catch (error) {
        console.error(error);
      }
    };

    const onSubmit = () => {
      newCoach();
    };
    const givenName = defineInputBinds("givenName", { validateOnChange: true });
    const middleName = defineInputBinds("middleName", {
      validateOnChange: true,
    });
    const surname = defineInputBinds("surname");
    const dateOfBirth = defineInputBinds("dateOfBirth", {
      validateOnChange: true,
    });
    const address = defineInputBinds("address", { validateOnChange: true });
    const email = defineInputBinds("email", { validateOnChange: true });
    const mobileNumber = defineInputBinds("mobileNumber", {
      validateOnChange: true,
    });
    const emergencyContact = defineInputBinds("emergencyContact", {
      validateOnChange: true,
    });

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
  compilerOptions: {
    delimiters: ["{", "}"],
  },
}).mount("#CoachRegistrationPage");
