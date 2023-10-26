import { createApp } from "vue";

createApp({
  setup() {
    const onSubmit = (event) => {
      const form = new FormData(event.target);
    };

    return {
      onSubmit,
    };
  },
}).mount("#ChangePassword");
