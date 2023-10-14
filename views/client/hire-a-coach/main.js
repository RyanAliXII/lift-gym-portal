import { createApp, onMounted, ref } from "vue";

import { Swiper, SwiperSlide } from "swiper/vue";
import Choices from "choices.js";

createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  components: {
    Swiper,
    SwiperSlide,
  },
  setup() {
    const coaches = ref([]);
    const coach = ref({
      givenName: "",
      surname: "",
      description: "",
      images: [],
    });
    const form = ref({
      coachId: 0,
      rateId: 0,
    });
    const slideTemplate = ref(null);
    const hireSelectElement = ref(null);
    const hireSelect = ref(null);
    const swiperElement = ref(null);
    const fetchCoaches = async () => {
      const response = await fetch("/clients/hire-a-coach", {
        headers: new Headers({ "Content-Type": "application/json" }),
      });

      if (response.status >= 400) return;
      const { data } = await response.json();
      coaches.value = data?.coaches ?? [];
    };
    const fetchCoachingRateByCoachId = async () => {
      try {
      } catch (error) {}
    };
    const preview = (c) => {
      c.description = c.description.replace("<script>", "");
      c.description = c.description.replace("</script>", "");
      coach.value = c;

      $("#profilePreviewModal").modal("show");
    };
    const initHire = (coachId) => {
      form.value.id = coachId;
      $("#hireModal").modal("show");
    };
    onMounted(() => {
      fetchCoaches();
      hireSelect.value = new Choices(hireSelectElement.value, {
        allowHTML: false,
      });
    });
    return {
      coaches,
      slideTemplate,
      swiperElement,
      preview,
      coach,
      initHire,
      hireSelectElement,
      publicUrl: window.publicUrl,
    };
  },
}).mount("#HireCoach");
// const swiper =
