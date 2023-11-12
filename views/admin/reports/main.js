const { createApp, ref, onMounted } = require("vue");
import VueDatePicker from "@vuepic/vue-datepicker";
import "@vuepic/vue-datepicker/dist/main.css";
createApp({
  components: {
    VueDatePicker,
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const dateRange = ref([]);

    const onSubmit = async () => {
      try {
        if (dateRange.value.length != 2) return;
        const start = dateRange.value[0];
        const end = dateRange.value[1];

        const response = await fetch("/app/reports", {
          method: "POST",
          body: JSON.stringify({
            dateRange: [start.toISOString(), end.toISOString()],
          }),
          headers: new Headers({
            "Content-Type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
        });
      } catch (err) {
        console.error(err);
      }
    };

    const setToMonthly = () => {
      const ONE_MONTH = 30; //30 days
      const endDate = new Date();
      const startDate = new Date(
        new Date().setDate(endDate.getDate() - ONE_MONTH)
      );
      dateRange.value = [startDate, endDate];
    };
    const setToWeekly = () => {
      const ONE_WEEK = 7; //7 days
      const endDate = new Date();
      const startDate = new Date(
        new Date().setDate(endDate.getDate() - ONE_WEEK)
      );
      dateRange.value = [startDate, endDate];
    };
    const setToAnnually = () => {
      const ONE_YEAR = 365; //365 days
      const endDate = new Date();
      const startDate = new Date(
        new Date().setDate(endDate.getDate() - ONE_YEAR)
      );
      dateRange.value = [startDate, endDate];
    };
    onMounted(() => {
      const endDate = new Date();
      const startDate = new Date(new Date().setDate(endDate.getDate() - 7));
      dateRange.value = [startDate, endDate];
    });
    return {
      dateRange,
      onSubmit,
      setToAnnually,
      setToWeekly,
      setToMonthly,
    };
  },
}).mount("#Reports");
