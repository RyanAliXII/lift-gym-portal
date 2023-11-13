import { createApp, onMounted, ref } from "vue";
import swal from "sweetalert2";
import { te } from "date-fns/locale";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const form = ref({ id: 0, meetingTime: "" });
    const errors = ref({});
    const appointments = ref([]);
    const fetchAppointments = async () => {
      try {
        const response = await fetch("/coaches/appointments", {
          headers: new Headers({
            "Content-Type": "application/json",
            "Cache-Control": "no-cache",
          }),
        });
        const { data } = await response.json();

        if (response.status === 200) {
          appointments.value = data?.appointments ?? [];
        }
      } catch (error) {
        console.error(error);
      }
    };
    const handleFormInput = (event) => {
      let value = event.target.value;
      let name = event.target.name;
      if (event.target.type === "number") {
        value = Number(value);
      }
      if (event.target.type === "datetime-local") {
        value = new Date(value).toISOString();
      }
      form.value[name] = value;
      delete errors.value[name];
    };
    const onSubmitApproval = async () => {
      try {
        errors.value = {};
        const response = await fetch(
          `/coaches/appointments/${form.value.id}/status?statusId=2`,
          {
            body: JSON.stringify(form.value),
            method: "PATCH",
            headers: new Headers({
              "Content-Type": "application/json",
              "X-CSRF-Token": window.csrf,
            }),
          }
        );
        const { data } = await response.json();
        if (response.status >= 400) {
          if (data?.errors) {
            errors.value = data?.errors;
          }
          return;
        }
        form.value = {
          id: 0,
          meetingTime: "",
        };
        $("#meetingDateModal").modal("hide");
        swal.fire(
          "Appointment Status Update",
          "Appointment status has been approved.",
          "success"
        );
        fetchAppointments();
      } catch (error) {
        console.error(error);
      }
    };

    const onSubmitPaid = async (id) => {
      try {
        errors.value = {};
        const response = await fetch(
          `/coaches/appointments/${id}/status?statusId=3`,
          {
            method: "PATCH",
            headers: new Headers({
              "Content-Type": "application/json",
              "X-CSRF-Token": window.csrf,
            }),
          }
        );

        swal.fire(
          "Appointment Status Update",
          "Appointment status has been mark as paid.",
          "success"
        );
        fetchAppointments();
      } catch (error) {
        console.error(error);
      }
    };
    const onSubmitCancellation = async (id, remarks) => {
      try {
        errors.value = {};
        const formData = new FormData();
        formData.append("remarks", remarks);
        const response = await fetch(
          `/coaches/appointments/${id}/status?statusId=4`,
          {
            body: formData,
            method: "PATCH",
            headers: new Headers({
              "X-CSRF-Token": window.csrf,
            }),
          }
        );
        if (response.status === 200) {
          swal.fire(
            "Appointment Status Update",
            "Appointment status has been cancelled.",
            "success"
          );
          fetchAppointments();
        }
      } catch (error) {
        console.error(error);
      }
    };
    const formatDate = (date) => {
      if (!date) return "";
      return new Date(date).toLocaleDateString("en-US", {
        month: "long",
        day: "2-digit",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
      });
    };
    const toMoney = (money) => {
      if (!money) return 0;
      return money.toLocaleString(undefined, {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
      });
    };

    const initApproval = (id) => {
      form.value.id = id;
      $("#meetingDateModal").modal("show");
    };

    const initCancellation = async (id) => {
      const { value: text, isConfirmed } = await swal.fire({
        input: "textarea",
        inputLabel: "Remarks",
        title: "Cancellation Remarks",
        confirmButtonText: "Submit",
        inputPlaceholder: "Enter the reason for cancellation.",
        inputAttributes: {
          "aria-label": "Enter the reason for cancellation.",
          maxlength: "150",
        },
        showCancelButton: true,
        confirmButtonColor: "#d9534f",
      });
      if (!isConfirmed) return;
      onSubmitCancellation(id, text);
    };
    const initMarkAsPaid = async (id) => {
      form.value = id;
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, this is already paid.",
        title: "Mark as Paid",
        text: "Are you that you want this appointment to mark as paid?",
        cancelButtonText: "This is not paid.",
        icon: "question",
      });
      if (result.isConfirmed) {
        onSubmitPaid(id);
      }
    };
    onMounted(() => {
      fetchAppointments();
    });
    const now = new Date().toISOString().slice(0, 16);
    return {
      appointments,
      toMoney,
      initApproval,
      now,
      handleFormInput,
      onSubmitApproval,
      errors,
      formatDate,
      initMarkAsPaid,
      initCancellation,
    };
  },
}).mount("#Appointments");
