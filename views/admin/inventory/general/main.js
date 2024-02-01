import { createApp, ref } from "vue";

createApp({
  setup() {
    const errors = ref({});
    const form = {
      name: "",
      brand: "",
      quantity: 0,
      unitOfMeasure: "",
      dateReceived: "",
      quantityThreshold: 0,
      costPrice: 0,
      image: "",
    };
    return {
      errors,
      form,
    };
  },
}).mount("#GeneralInventoryPage");
