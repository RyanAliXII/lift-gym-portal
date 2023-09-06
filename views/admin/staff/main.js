import { useForm } from "vee-validate";
import { createApp } from "vue";
import swal from "sweetalert2";
createApp({
  setup() {
    const {
      defineInputBinds,
      values: form,
      errors,
    } = useForm({
      initialValues: {
        givenName: "",
        middleName: "",
        surname: "",
        email: "",
      },
    });

    const givenName = defineInputBinds("givenName", {
      validateOnChange: true,
    });
    const middleName = defineInputBinds("middleName", {
      validateOnChange: true,
    });
    const surname = defineInputBinds("surname", {
      validateOnChange: true,
    });
    const email = defineInputBinds("email", {
      validateOnChange: true,
    });

    const onSubmitNewStaff = async () => {
      try {
        const response = await fetch("/app/staffs", {
          method: "POST",
          body: JSON.stringify(form),
          headers: new Headers({
            "content-type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
        });
        if (response.status === 200) {
          swal.fire(
            "New Staff",
            "Staff has been succcessfully added.",
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
      email,
      errors,
      onSubmitNewStaff,
    };
  },
}).mount("#StaffPage");
