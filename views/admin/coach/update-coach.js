import { useForm } from "vee-validate";
import { createApp } from "vue";
import swal from "sweetalert2";
import { object } from "yup";
createApp({
  setup() {
    const { errors, defineInputBinds, values, setErrors } = useForm({
      initialValues: window.coach ?? {
        id: 0,
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

    const onSubmit = async () => {
      try {
        const response = await fetch(`/app/coaches/${window.coach.id ?? 0}`, {
          headers: new Headers({
            "content-type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
          method: "PUT",
          body: JSON.stringify(values),
        });

        const { data } = await response.json();
        if (response.status === 200) {
          swal.fire(
            "Coach Profile Update",
            `Coach profile has been updated.`,
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
    const initializeResetPassword = (id) => {
      swal
        .fire({
          title: "Reset Password",
          text: "Are you sure you want to reset this account password?",
          icon: "warning",
          showCancelButton: true,
          cancelButtonColor: "#3085d6",
          confirmButtonColor: "#d33",
          confirmButtonText: "Yes, reset it",
          cancelButtonText: "I don't want to reset the password",
        })
        .then((result) => {
          if (result.isConfirmed) {
            resetPassword(id);
          }
        });
    };
    const resetPassword = async (id) => {
      try {
        const response = await fetch(`/app/coaches/${id}/password`, {
          method: "PATCH",
          headers: new Headers({
            "content-type": "application/json",
            "X-CSRF-Token": csrf,
          }),
        });

        const { data } = await response.json();
        if (response.status === 200) {
          swal.fire(
            "Password Reset",
            `Password has been reset. <br> The new password for the account is <strong>${data?.password}</strong>.<br>Please keep the password this will be the only time it will be shown.`,
            "success"
          );
        }
      } catch (error) {
        console.error(error);
      }
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
      initializeResetPassword,
    };
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
}).mount("#CoachRegistrationPage");
