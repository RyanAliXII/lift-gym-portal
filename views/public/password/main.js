import { tr } from "date-fns/locale";
import { createApp, onMounted, ref } from "vue";

createApp({
  setup() {
    const displaySuccessMessage = ref(false);
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
        if (response.status === 200) {
          displaySuccessMessage.value = true;
        }
      } catch (error) {
        console.log(error);
      }
    };

    return {
      onSubmit,
      displaySuccessMessage,
    };
  },
}).mount("#ResetPassword");
