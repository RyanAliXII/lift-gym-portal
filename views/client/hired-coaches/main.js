import { createApp, onMounted, ref } from "vue";

createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const hiredCoaches = ref([]);
    const fetchHiredCoaches = async () => {
      try {
        const response = await fetch("/clients/hired-coaches", {
          headers: new Headers({ "Content-Type": "application/json" }),
        });
        const { data } = await response.json();

        if (response.status === 200) {
          hiredCoaches.value = data?.hiredCoaches ?? [];
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
      fetchHiredCoaches();
    });
    return {
      hiredCoaches,
      toMoney,
    };
  },
}).mount("#HiredCoaches");
