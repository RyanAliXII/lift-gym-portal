import { createApp, onMounted, ref } from "vue";

createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const reservations = ref([]);
    const fetchReservations = async () => {
      const response = await fetch("/app/reservations", {
        headers: new Headers({ "Content-Type": "application/json" }),
      });
      const { data } = await response.json();
      if (response.status === 200) {
        reservations.value = data?.reservations ?? [];
      }
    };
    const formatDate = (date) => {
      if (!date) return "";
      return new Date(date).toLocaleDateString("en-US", {
        month: "long",
        day: "2-digit",
        year: "numeric",
      });
    };
    const handleDateSelect = (event) => {
      const id = event.target.value;
      if (id == 0) {
        fetchReservations();
        return;
      }
      fetchReservationsDateSlot(id);
    };
    const fetchReservationsDateSlot = async (id) => {
      const response = await fetch(`/app/reservations/date-slots/${id}`, {
        headers: new Headers({ "Content-Type": "application/json" }),
      });
      const { data } = await response.json();
      if (response.status === 200) {
        reservations.value = data?.reservations ?? [];
      }
    };
    onMounted(() => {
      fetchReservations();
    });
    return {
      reservations,
      formatDate,
      handleDateSelect,
    };
  },
}).mount("#ReservationPage");
