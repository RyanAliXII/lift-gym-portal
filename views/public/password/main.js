import { createApp, onMounted, ref } from "vue";

createApp({
  setup() {
    const errors = ref({});

    const onSubmit = async (event) => {
      const form = new FormData(event.target);
      try {
        const response = await fetch("/app/reset-password", {
          method: "POST",
          body: form,
          headers: new Headers({
            "X-CSRF-Token": window.csrf,
          }),
        });
      } catch (error) {
        console.log(error);
      }
    };
    onMounted(() => {
      console.log("test");
    });
    return {
      onSubmit,
    };
  },
}).mount("#ResetPassword");
