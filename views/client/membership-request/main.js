import { createApp, onMounted, ref } from "vue";
import Choices from "choices.js";
createApp({
  setup() {
    let planSelectElement = ref(null);
    let planSelect = null;
    onMounted(() => {
      planSelect = new Choices(planSelectElement.value, {});
    });
    return {
      planSelectElement,
    };
  },
}).mount("#MembershipRequest");
