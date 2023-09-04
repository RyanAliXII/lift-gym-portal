import { createApp, onMounted, ref } from "vue";
import Choices from "choices.js";

createApp({
  setup() {
    const packageSelectElement = ref(null);
    let packageSelect = null;

    const errorMessage = ref(undefined);

    const fetchPackages = async () => {
      try {
        const response = await fetch("/clients/packages");
        const { data } = await response.json();
        return data?.packages ?? [];
      } catch {
        return [];
      }
    };
    const init = async () => {
      packageSelect = new Choices(packageSelectElement.value, {});
      const packages = await fetchPackages();
      const packageOptions = packages.map((p) => ({
        value: p.id,
        label: p.description + " - " + p.price + " PHP",
        id: p.id,
        customProperties: p,
      }));
      packageSelect.setChoices(packageOptions);
    };
    onMounted(() => {
      init();
    });

    return {
      packageSelectElement,
      errorMessage,
    };
  },
}).mount("#PackageRequest");
