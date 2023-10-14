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
        fetchRates();
      } catch (error) {
        console.error(error);
      }
    };

    const onSubmitUpdateRate = async () => {
      try {
        errors.value = {};
        const response = await fetch(`/coaches/rates/${form.value.id}`, {
          method: "PUT",
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
        swal.fire("Rate Update", "Rate has been updated.", "success");
        $("#editRateModal").modal("hide");
        fetchRates();
      } catch (error) {
        console.error(error);
      }
    };

    const deleteRate = async (rateId) => {
      const response = await fetch(`/coaches/rates/${rateId}`, {
        method: "DELETE",
        headers: new Headers({
          "Content-Type": "application/json",
          "X-CSRF-Token": window.csrf,
        }),
      });
      if (response.status >= 400) return;
      swal.fire(
        "Coaching Rate Deletion",
        "Coaching rate has been deleted.",
        "success"
      );
      fetchRates();
    };
    const initDelete = async (id) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, delete it",
        title: "Delete Coaching Rate",
        text: "Are you sure you want to delete coaching rate?",
        confirmButtonColor: "#d9534f",
        cancelButtonText: "I don't want to delete this coaching rate",
        icon: "warning",
      });
      if (result.isConfirmed) {
        deleteRate(id);
      }
    };
    onMounted(() => {
      fetchRates();
      $("#newRateModal").on("hidden.bs.modal", () => {
        form.value = { ...initialForm };
        errors.value = {};
      });
      $("#editRateModal").on("hidden.bs.modal", () => {
        form.value = { ...initialForm };
        errors.value = {};
      });
    });
    const initEdit = (rate) => {
      form.value = rate;
      $("#editRateModal").modal("show");
    };

    return {
      form,
      errors,
      handleFormInput,
      rates,
      initEdit,
      onSubmitNewRate,
      onSubmitUpdateRate,
      initDelete,
    };
  },
}).mount("#CoachingRate");
