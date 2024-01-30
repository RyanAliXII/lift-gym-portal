import { createApp, onMounted, ref } from "vue";
import swal from "sweetalert2";
import { parse, format } from "date-fns";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const form = ref({
      id: 0,
      date: "",
      time: "",
    });
    const schedules = ref([]);
    const errors = ref({});
    const fetchSchedules = async () => {
      try {
        const response = await fetch("/coaches/schedules", {
          headers: new Headers({
            "Content-Type": "application/json",
            "Cache-Control": "no-cache",
          }),
        });
        const { data } = await response.json();
        schedules.value = data?.schedules ?? [];
      } catch (error) {
        console.error(error);
      }
    };
    const resetForm = () => {
      form.value = {
        date: "",
        id: 0,
        time: "",
      };
      errors.value = {};
    };
    onMounted(() => {
      fetchSchedules();
      $("#addSchedModal").on("hidden.bs.modal", () => {
        resetForm();
      });

      $("#editSchedModal").on("hidden.bs.modal", () => {
        resetForm();
      });
    });
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
        fetchSchedules();

        $("#addSchedModal").modal("hide");
      } catch (error) {
        console.error(error);
      }
    };

    const deleteSchedule = async (id) => {
      try {
        await fetch(`/coaches/schedules/${id}`, {
          method: "DELETE",
          headers: new Headers({
            "Content-type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
        });
        swal.fire("Schedule", "Schedule has been deleted.", "success");
        fetchSchedules();
      } catch (error) {
        console.error(error);
      }
    };
    const onSubmitUpdate = async () => {
      try {
        errors.value = {};
        const response = await fetch(`/coaches/schedules/${form.value.id}`, {
          method: "PUT",
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
        swal.fire("Schedule", "Schedule has been updated.", "success");
        form.value = {
          date: "",
          time: "",
        };
        fetchSchedules();
        $("#editSchedModal").modal("hide");
      } catch (error) {
        console.error(error);
      }
    };
    const initEditSched = (data) => {
      form.value = { ...data };
      $("#editSchedModal").modal("show");
    };

    const initDelete = async (id) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, delete it",
        title: "Delete Schedule?",
        text: "Are you sure you want to delete schedule?",
        confirmButtonColor: "#d9534f",
        cancelButtonText: "I don't want to delete this schedule",
        icon: "warning",
      });
      if (result.isConfirmed) {
        deleteSchedule(id);
      }
    };
    return {
      form,
      onSubmit,
      errors,
      schedules,
      toReadableDate,
      initEditSched,
      onSubmitUpdate,
      initDelete,
      to12HR,
    };
  },
}).mount("#CoachingSchedule");
