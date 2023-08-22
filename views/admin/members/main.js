import { createApp, ref, onMounted } from "vue";
createApp({
  setup() {
    const message = "Sample message from this app";

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
}).mount("#MembersPage");
