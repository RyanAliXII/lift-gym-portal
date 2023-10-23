import { createApp, onMounted, ref } from "vue";

createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const appointments = ref([]);
    const fetchAppoinments = async () => {
      try {
        const response = await fetch("/coaches/appointments", {
          headers: new Headers({ "Content-Type": "application/json" }),
        });
        const { data } = await response.json();

        if (response.status === 200) {
          appointments.value = data?.appointments ?? [];
        }
      } catch (error) {
        console.error(error);
      }
    };
    const toMoney = (money) => {
      if (!money) return 0;
      return money.toLocaleString(undefined, {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
      });
    };
    onMounted(() => {
      fetchAppoinments();
    });
    return {
      appointments,
      toMoney,
    };
  },
}).mount("#Appointments");
