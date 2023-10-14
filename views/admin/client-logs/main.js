import Choices from "choices.js";
import { createApp, onMounted, ref } from "vue";
import { useDebounce, useDebounceFn } from "@vueuse/core";
createApp({
  setup() {
    const logClientSelectElement = ref(null);
    const logClientSelect = ref(null);

    const search = useDebounceFn((query) => {}, 500);

    onMounted(() => {
      logClientSelect.value = new Choices(logClientSelectElement.value, {
        allowHTML: false,
        placeholder: "Seach Client",
      });

      logClientSelect.value.passedElement.element.addEventListener(
        "search",
        (event) => {
          search(event.detail.value);
        }
      );
    });

    return {
      logClientSelectElement,
    };
  },
}).mount("#ClientLog");
