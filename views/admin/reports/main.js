const { createApp, ref, onMounted } = require("vue");
import VueDatePicker from "@vuepic/vue-datepicker";
import "@vuepic/vue-datepicker/dist/main.css";
createApp({
  components: {
    VueDatePicker,
  },

  setup() {
    const dateRange = ref([]);
    const today = new Date();
    const onSubmit = () => {
      if (dateRange.value.length != 2) return;
      const start = dateRange.value[0];
      const end = dateRange.value[1];
    };

    const setToMonthly = () => {
      const ONE_MONTH = 30;
      const endDate = new Date();
      const startDate = new Date(new Date().setDate(endDate.getDate() - 30));
    };
    const setToWeekly = () => {};
    onMounted(() => {
      const endDate = new Date();
      const startDate = new Date(new Date().setDate(endDate.getDate() - 7));
      dateRange.value = [startDate, endDate];
    });
    return {
      dateRange,
      onSubmit,
    };
  },
}).mount("#Reports");
