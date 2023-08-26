import { createApp, ref } from "vue";

createApp({
  setup() {
    const form = ref({
      email: "",
      password: "",
    });

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
      console.log(data);
    };
    return {
      form,
      onSubmit,
    };
  },
}).mount("#ClientLoginPage");
