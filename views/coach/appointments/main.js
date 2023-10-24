import { createApp, onMounted, ref } from "vue";
import swal from "sweetalert2";
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
          headers: new Headers({ "Content-Type": "application/json" }),
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
    };
  },
}).mount("#Appointments");
