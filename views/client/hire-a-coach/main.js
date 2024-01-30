import { createApp, onMounted, ref } from "vue";
import { Swiper, SwiperSlide } from "swiper/vue";
import Choices from "choices.js";
import swal from "sweetalert2";
import { format, parse } from "date-fns";
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
      meetingTime: "",
      scheduleId: 0,
    });
    const today = new Date();
    today.setDate(today.getDate() + 1);
    const minDate = today.toISOString().slice(0, -8);

    const slideTemplate = ref(null);
    const hireSelectElement = ref(null);
    const scheduleSelectElement = ref(null);
    const hireSelect = ref(null);
    const scheduleSElect = ref(null);
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
      fetchCoachScheds(coachId);
      $("#hireModal").modal("show");
    };

    const onSubmit = async () => {
      try {
        errors.value = {};
        const response = await fetch("/clients/hire-a-coach", {
          body: JSON.stringify({
            ...form.value,
          }),
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
          scheduleId: 0,
        };
        $("#hireModal").modal("hide");
        swal.fire("Hire Coach", "Coach has been hired.", "success");
        errors.value = {};
      } catch (error) {
        console.error(error);
      }
    };

    const toReadableDate = (d) => {
      if (!d) return "";
      const dt = new Date(d);
      try {
        return dt.toLocaleDateString(undefined, {
          month: "long",
          year: "numeric",
          day: "2-digit",
        });
      } catch (error) {
        return "";
      }
    };

    const to12HR = (timeStr) => {
      if (!timeStr) return "";
      try {
        const parsedTime = parse(timeStr, "HH:mm:ss", new Date());
        const formattedTime = format(parsedTime, "h:mm a");
        return formattedTime;
      } catch (error) {
        console.error(error);
        return "";
      }
    };
    const fetchCoachScheds = async (id) => {
      try {
        scheduleSElect.value.clearStore();
        const response = await fetch(`/clients/coaches/${id}/schedules`);
        const { data } = await response.json();
        const schedSelectValues = (data?.schedules ?? []).map((sched) => ({
          value: sched.id,
          label: `${toReadableDate(sched.date)} ${to12HR(sched.time)}`,
        }));
        scheduleSElect.value.setChoices(
          schedSelectValues,
          "value",
          "label",
          true
        );
      } catch (error) {
        console.error(error);
      }
    };
    onMounted(() => {
      fetchCoaches();
      hireSelect.value = new Choices(hireSelectElement.value, {
        allowHTML: false,
      });
      scheduleSElect.value = new Choices(scheduleSelectElement.value, {
        allowHTML: false,
      });

      hireSelect.value.passedElement.element.addEventListener(
        "change",
        (event) => {
          form.value.rateId = event.detail.value;
        }
      );
      scheduleSElect.value.passedElement.element.addEventListener(
        "change",
        (event) => {
          form.value.scheduleId = event.detail.value;
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
      form,
      scheduleSelectElement,
      minDate,
    };
  },
}).mount("#HireCoach");
// const swiper =
