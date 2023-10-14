import { createApp, onMounted, ref } from "vue";
import Swiper from "swiper";
import "swiper/css";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const coaches = ref([]);
    const coach = ref({
      givenName: "",
      surname: "",
      description: "",
    });
    const swiperElement = ref(null);
    const swiper = ref(null);
    const fetchCoaches = async () => {
      const response = await fetch("/clients/hire-a-coach", {
        headers: new Headers({ "Content-Type": "application/json" }),
      });

      if (response.status >= 400) return;
      const { data } = await response.json();
      coaches.value = data?.coaches ?? [];
    };
    const preview = (coach) => {
      console.log(coach);
      $("#profilePreviewModal").modal("show");
    };

    onMounted(() => {
      fetchCoaches();
      swiper.value = new Swiper(swiperElement.value, {
        // Optional parameters
        direction: "horizontal",
        loop: true,

        // If we need pagination
        pagination: {
          el: ".swiper-pagination",
        },

        // Navigation arrows
        navigation: {
          nextEl: ".swiper-button-next",
          prevEl: ".swiper-button-prev",
        },
      });
    });
    return {
      coaches,
      swiperElement,
      preview,
    };
  },
}).mount("#HireCoach");
// const swiper =
