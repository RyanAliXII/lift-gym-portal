import Swal from "sweetalert2";
import { createApp, ref } from "vue";

createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const form = ref({
      name: "",
      message: "",
      email: "",
    });

    const onSubmit = async () => {
      const response = await fetch("/contact-us", {
        method: "POST",
        body: JSON.stringify(form.value),
        headers: new Headers({
          "Content-Type": "application/json",
          "X-CSRF-Token": window.csrf,
        }),
      });
      if (response.status >= 400) {
        const { data } = await response.json();
        if (data?.errors) {
          Swal.fire("Error", "Please fill up required fields.", "error");
          return;
        }
        Swal.fire("Error", "Unknown error occured.", "error");
        return;
      }
      form.value = {
        name: "",
        message: "",
        email: "",
      };
      Swal.fire("Message Sent", "Your message has been sent.", "success");
    };

    return {
      form,

      onSubmit,
    };
  },
}).mount("#ContactPage");
