import { createApp, onMounted, ref } from "vue";
import swal from "sweetalert2";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const initialForm = {
      id: 0,
      description: "",
      price: 0,
    };
    const form = ref({ ...initialForm });
    const handleFormInput = (event) => {
      let value = event.target.value;
      let name = event.target.name;
      if (event.target.type === "number") {
        value = Number(value);
      }
      form.value[name] = value;
      delete errors.value[name];
    };
    const errors = ref({});
    const rates = ref([]);
    const fetchRates = async () => {
      try {
        const response = await fetch("/coaches/rates", {
          headers: new Headers({ "Content-Type": "application/json" }),
        });
        const { data } = await response.json();
        if (response.status >= 500) return;
        rates.value = data?.rates ?? [];
      } catch (error) {
        console.error(error);
      }
    };
    const onSubmitNewRate = async () => {
      try {
        errors.value = {};
        const response = await fetch("/coaches/rates", {
          method: "POST",
          body: JSON.stringify(form.value),
          headers: new Headers({
            "Content-Type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
        });
        const { data } = await response.json();

        if (response.status >= 400) {
          if (data?.errors) {
            errors.value = data?.errors ?? {};
          }
          return;
        }
        swal.fire("New Rate", "New rate has been created.", "success");
        $("#newRateModal").modal("hide");
      } catch (error) {
        console.error(error);
      }
    };
    onMounted(() => {
      fetchRates();
    });

    return {
      form,
      errors,
      handleFormInput,
      rates,
      onSubmitNewRate,
    };
  },
}).mount("#CoachingRate");
