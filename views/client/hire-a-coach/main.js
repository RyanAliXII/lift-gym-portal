import { createApp, onMounted, ref } from "vue";
import Swiper from "swiper";
import "swiper/css";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const coaches = ref([]);
    const fetchCoaches = async () => {
      const response = await fetch("/clients/hire-a-coach", {
        headers: new Headers({ "Content-Type": "application/json" }),
      });

      if (response.status >= 400) return;

      const { data } = await response.json();
      coaches.value = data?.coaches ?? [];
    };

    onMounted(() => {
      fetchCoaches();
    });
    return {
      coaches,
    };
  },
}).mount("#HireCoach");
