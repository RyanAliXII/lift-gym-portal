import { createApp, ref, onMounted } from "vue";
createApp({
  setup() {
    const message = "Subs page";
    onMounted(() => {
      console.log("App Mounted");
    });
    return {
      message,
    };
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
}).mount("#SubscribePage");
