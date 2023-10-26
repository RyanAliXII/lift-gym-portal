import { createApp, ref } from "vue";
import swal from "sweetalert2";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const errors = ref({});
    const onSubmit = async (event) => {
      errors.value = {};
      const form = new FormData(event.target);
      try {
        const response = await fetch("/app/profile/password", {
          method: "PATCH",
          body: form,
          headers: new Headers({
            "X-CSRF-Token": window.csrf,
          }),
        });
        const { data } = await response.json();
        if (response.status >= 400) {
          if (data?.errors) {
            errors.value = data?.errors;
          }
          return;
        }
        swal.fire(
          "Change Password",
          "Your password has been changed.",
          "success"
        );
        event.target.reset();
      } catch (err) {
        console.error(err);
      }
    };
    return {
      onSubmit,
      errors,
    };
  },
}).mount("#ProfilePage");
