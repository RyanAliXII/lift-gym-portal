import swal from "sweetalert2";
import { Form } from "vee-validate";
import { createApp, onMounted, ref } from "vue";

createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const initialValues = {
      startTime: "",
      endTime: "",
      maxCapacity: 20,
    };
    const form = ref({ ...initialValues });
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
        const response = await fetch("/app/time-slots", {
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
            errors.value = data.errors;
          }
          return;
        }
        $("#newSlotModal").modal("hide");
        form.value = { ...initialValues };
        swal.fire("New Time Slot", "Time slot has been created.", "success");
      } catch (error) {
        console.error(error);
      }
    };
    onMounted(() => {});
    return {
      form,
      handleFormInput,
      onSubmit,
      errors,
    };
  },
}).mount("#TimeSlot");
