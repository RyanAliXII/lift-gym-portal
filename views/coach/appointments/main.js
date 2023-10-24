import { createApp, onMounted, ref } from "vue";

createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const form = ref({ id: 0, meetingDate: "" });
    const errors = ref({});
    const appointments = ref([]);
    const fetchAppoinments = async () => {
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
      console.log(event);
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
    const toMoney = (money) => {
      if (!money) return 0;
      return money.toLocaleString(undefined, {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
      });
    };
    const onSubmitMeetingDate = () => {};
    const initApproval = (id) => {
      form.value.id = id;
      $("#meetingDateModal").modal("show");
    };
    onMounted(() => {
      fetchAppoinments();
    });
    const now = new Date().toISOString().slice(0, 16);
    return {
      appointments,
      toMoney,
      initApproval,
      now,
      handleFormInput,
    };
  },
}).mount("#Appointments");
