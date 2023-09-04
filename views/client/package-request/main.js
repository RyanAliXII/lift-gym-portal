import { createApp, onMounted, ref } from "vue";
import Choices from "choices.js";

createApp({
  setup() {
    const packageSelectElement = ref(null);
    let packageSelect = null;
    const errorMessage = ref(undefined);
    onMounted(() => {
      packageSelect = new Choices(packageSelectElement.value, {});
    });

    return {
      packageSelectElement,
      errorMessage,
    };
  },
}).mount("#PackageRequest");
