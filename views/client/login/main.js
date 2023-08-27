import { createApp, ref } from "vue";

createApp({
  setup() {
    const form = ref({
      email: "",
      password: "",
    });
    const message = ref("");
    const onSubmit = async () => {
      const response = await fetch("/clients/login", {
        body: JSON.stringify(form.value),
        method: "POST",
        headers: new Headers({
          "Content-Type": "application/json",
          "X-CSRF-Token": window.csrf,
        }),
      });
      const data = await response.json();
      if (response.status === 200) {
        location.replace("/clients/dashboard");
      }
      message.value = data?.message;
    };

    return {
      form,
      message,
      onSubmit,
    };
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
}).mount("#ClientLoginPage");
