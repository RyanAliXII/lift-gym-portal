import { createApp, ref } from "vue";

createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const errors = ref({});
    const showPassword = ref(false);
    const formValue = ref({
      password: "",
      confirmPassword: "",
    });
    const displaySuccessMessage = ref(false);
    const onSubmit = async (event) => {
      const form = new FormData(event.target);

      const searchParams = new URLSearchParams(window.location.search);
      const key = searchParams.get("key");
      form.set("publicKey", key);
      try {
        errors.value = {};
        const response = await fetch("/change-password", {
          method: "POST",
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
        displaySuccessMessage.value = true;
      } catch (error) {
        console.error(error);
      }
    };

    return {
      onSubmit,
      errors,
      showPassword,
      formValue,
      displaySuccessMessage,
    };
  },
}).mount("#ChangePassword");
