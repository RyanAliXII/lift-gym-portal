import { format } from "date-fns";
import { createApp, ref } from "vue";

createApp({
  setup() {
    const today = format(new Date(), "yyyy-MM-dd");
    const form = ref({
      to: today,
      from: today,
    });
    const errors = ref({});
    return {
      form,
    };
  },
}).mount("#DateSlot");
