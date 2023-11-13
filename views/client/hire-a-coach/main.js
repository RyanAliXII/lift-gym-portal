import { createApp, onMounted, ref } from "vue";
import { Swiper, SwiperSlide } from "swiper/vue";
import Choices from "choices.js";
import swal from "sweetalert2";
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
    const errors = ref({
      rateId: undefined,
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
        headers: new Headers({
          "Content-Type": "application/json",
          "Cache-Control": "no-cache",
        }),
      });

      if (response.status >= 400) return;
      const { data } = await response.json();
      coaches.value = data?.coaches ?? [];
    };
    const fetchCoachingRatesByCoachId = async (coachId) => {
      try {
        const response = await fetch(`/clients/coaches/${coachId}/rates`);

        if (response.status >= 400) return;
        const { data } = await response.json();
        const rateSelectValues = (data?.rates ?? []).map((rate) => ({
          value: rate.id,
          label: rate.description,
        }));
        hireSelect.value.setChoices(rateSelectValues, "value", "label", true);
      } catch (error) {}
    };
    const preview = (c) => {
      c.description = c.description.replace("<script>", "");
      c.description = c.description.replace("</script>", "");
      coach.value = c;

      $("#profilePreviewModal").modal("show");
    };
    const initHire = (coachId) => {
      form.value.coachId = coachId;
      fetchCoachingRatesByCoachId(coachId);
      $("#hireModal").modal("show");
    };

    const onSubmit = async () => {
      try {
        const response = await fetch("/clients/hire-a-coach", {
          body: JSON.stringify(form.value),
          method: "POST",
          headers: new Headers({
            "Content-Type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
        });
        const { data } = await response.json();
        if (response.status >= 400) {
          if (data?.errors) {
            errors.value = data?.errors;
            return;
          }
          return;
        }
        form.value = {
          coachId: 0,
          rateId: 0,
        };
        $("#hireModal").modal("hide");
        swal.fire("Hire Coach", "Coach has been hired.", "success");
      } catch (error) {
        console.error(error);
      }
    };
    onMounted(() => {
      fetchCoaches();
      hireSelect.value = new Choices(hireSelectElement.value, {
        allowHTML: false,
      });
      hireSelect.value.passedElement.element.addEventListener(
        "change",
        (event) => {
          form.value.rateId = event.detail.value;
        }
      );
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
      onSubmit,
      errors,
    };
  },
}).mount("#HireCoach");
// const swiper =
