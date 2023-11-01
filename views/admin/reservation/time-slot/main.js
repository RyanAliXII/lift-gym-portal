import swal from "sweetalert2";
import { computed, createApp, onMounted, ref } from "vue";
import { parse, format } from "date-fns";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const initialValues = {
      id: 0,
      startTime: "",
      endTime: "",
      maxCapacity: 20,
    };
    const form = ref({ ...initialValues });
    const timeSlots = ref([]);
    const timeSlotSelections = ref([]);
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
        fetchTimeSlots();
        swal.fire("New Time Slot", "Time slot has been created.", "success");
      } catch (error) {
        console.error(error);
      }
    };

    const onSubmitUpdate = async () => {
      try {
        errors.value = {};
        const response = await fetch(`/app/time-slots/${form.value.id}`, {
          body: JSON.stringify(form.value),
          method: "PUT",
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
        $("#editSlotModal").modal("hide");
        form.value = { ...initialValues };
        fetchTimeSlots();
        swal.fire("Update Time Slot", "Time slot has been updated.", "success");
      } catch (error) {
        console.error(error);
      }
    };
    const fetchTimeSlots = async () => {
      try {
        const response = await fetch("/app/time-slots", {
          headers: new Headers({ "Content-Type": "application/json" }),
        });
        const { data } = await response.json();

        if (response.status === 200) {
          timeSlots.value = data?.slots ?? [];
          timeSlotSelections.value = data?.slotSelections ?? [];
        }
      } catch (error) {
        console.error(error);
      }
    };
    const formatTime = (time) => {
      if (!time) return "";
      try {
        const parsedTime = parse(time, "HH:mm:ss", new Date());
        const formattedTime = format(parsedTime, "h:mm a");
        return formattedTime;
      } catch (error) {
        console.error(error);
        return "";
      }
    };

    const deleteTimeSlot = async (id) => {
      const response = await fetch(`/app/time-slots/${id}`, {
        method: "DELETE",
        headers: new Headers({
          "X-CSRF-Token": window.csrf,
        }),
      });
      if (response.status === 200) {
        swal.fire("Delete Time Slot", "Time slot has been deleted.", "success");
        fetchTimeSlots();
      }
    };
    const initDelete = async (id) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, delete it",
        title: "Delete Time Slot",
        text: "Are you sure you want to delete slot",
        confirmButtonColor: "#d9534f",
        cancelButtonText: "I don't want to delete this slot",
        icon: "warning",
      });
      if (result.isConfirmed) {
        deleteTimeSlot(id);
      }
    };
    const initEdit = (slot) => {
      try {
        const startTime = parse(slot.startTime, "HH:mm:ss", new Date());
        const endTime = parse(slot.endTime, "HH:mm:ss", new Date());
        const formattedStart = format(startTime, "HH:mm");
        const formattedEnd = format(endTime, "HH:mm");
        form.value = {
          id: slot.id,
          startTime: formattedStart,
          endTime: formattedEnd,
          maxCapacity: slot.maxCapacity,
        };
        $("#editSlotModal").modal("show");
      } catch (err) {
        console.error(err);
      }
    };
    onMounted(() => {
      fetchTimeSlots();
      $("#newSlotModal").on("hidden.bs.modal", () => {
        errors.value = {};
        form.value = { ...initialValues };
      });
      $("#editSlotModal").on("hidden.bs.modal", () => {
        errors.value = {};
        form.value = { ...initialValues };
      });
    });
    return {
      form,
      handleFormInput,
      onSubmit,
      formatTime,
      timeSlots,
      errors,
      initDelete,
      initEdit,
      onSubmitUpdate,
      timeSlotSelections,
    };
  },
}).mount("#TimeSlot");
