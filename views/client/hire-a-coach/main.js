import { createApp, onMounted, ref } from "vue";

import { Swiper, SwiperSlide } from "swiper/vue";

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
    const slideTemplate = ref(null);
    const swiperElement = ref(null);
    const fetchCoaches = async () => {
      const response = await fetch("/clients/hire-a-coach", {
        headers: new Headers({ "Content-Type": "application/json" }),
      });

      if (response.status >= 400) return;
      const { data } = await response.json();
      coaches.value = data?.coaches ?? [];
    };
    const preview = (c) => {
      c.description = c.description.replace("<script>", "");
      c.description = c.description.replace("</script>", "");
      coach.value = c;

      $("#profilePreviewModal").modal("show");
    };

    onMounted(() => {
      fetchCoaches();
    });
    return {
      coaches,
      slideTemplate,
      swiperElement,
      preview,
      coach,
      publicUrl: window.publicUrl,
    };
  },
}).mount("#HireCoach");
// const swiper =
