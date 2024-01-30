import { createApp, onMounted, ref } from "vue";
import swal from "sweetalert2";

createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const form = ref({
      date: "",
      time: "",
    });

    const errors = ref({});
    const onSubmit = async () => {
      try {
        errors.value = {};
        const response = await fetch("/coaches/schedules", {
          method: "POST",
          body: JSON.stringify(form.value),
          headers: new Headers({
            "Content-type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
        });

        const { data } = await response.json();
        if (response.status >= 400) {
          if (data?.errors) {
            errors.value = data?.errors;
          }
          return;
        }
        swal.fire("Schedule", "Schedule has been created.", "success");
        form.value = {
          date: "",
          time: "",
        };
        $("#addSchedModal").modal("hide");
      } catch (error) {
        console.error(error);
      }
    };
    return {
      form,
      onSubmit,
      errors,
    };
  },
}).mount("#CoachingSchedule");
