import { format } from "date-fns";
import { createApp, onMounted, ref } from "vue";
import swal from "sweetalert2";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const today = format(new Date(), "yyyy-MM-dd");
    const initialFormValue = {
      to: today,
      from: today,
    };
    const form = ref({
      ...initialFormValue,
    });
    const slots = ref([]);
    const errors = ref({});
    const handleFormInput = (event) => {
      let value = event.target.value;
      let name = event.target.name;
      if (event.target.type === "number") {
        value = Number(value);
      }
      form.value[name] = value;
      delete errors.value[name];
    };
    const onSubmit = async () => {
      try {
        errors.value = {};
        const response = await fetch("/app/date-slots", {
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
            errors.value = data.errors;
          }
          return;
        }
        form.value = { ...initialFormValue };
        $("#newSlotModal").modal("hide");
        swal.fire("New Slot", "Date slot/s has been added.", "success");
      } catch (error) {
        console.error(error);
      }
    };
    const fetchSlots = async () => {
      try {
        const response = await fetch("/app/date-slots", {
          headers: new Headers({ "Content-Type": "application/json" }),
        });
        const { data } = await response.json();
        slots.value = data?.slots ?? [];
      } catch (error) {
        console.error(error);
      }
    };

    const formatDate = (date) => {
      if (!date) return "No Date";
      if (date.length === 0) return "No Date";
      return new Date(date).toLocaleDateString("en-US", {
        month: "long",
        day: "2-digit",
        year: "numeric",
      });
    };
    onMounted(() => {
      fetchSlots();
    });
    return {
      form,
      handleFormInput,
      onSubmit,
      errors,
      today,
      formatDate,
      slots,
    };
  },
}).mount("#DateSlot");
