import { createApp, onMounted, ref } from "vue";
import { Calendar } from "fullcalendar";
import interactionPlugin from "@fullcalendar/interaction";
import { format } from "date-fns";
createApp({
  setup() {
    const reservationCalendarElement = ref(null);
    const reservationCalendar = ref(null);
    const dateSlots = ref([]);
    const timeSlots = ref([]);
    const selectedDate = ref("");
    const fetchDateSlots = async () => {
      try {
        const response = await fetch("/clients/reservations/date-slots");
        const { data } = await response.json();
        dateSlots.value = data?.slots ?? [];
      } catch (error) {
        console.error(error);
      }
    };
    const fetchTimeSlotsBasedOnDateSlot = async (id) => {
      try {
        const response = await fetch(
          `/clients/reservations/date-slots/${id}/time-slots`
        );
        const { data } = await response.json();
        timeSlots.value = data?.slots ?? [];
      } catch (error) {
        console.error(error);
      }
    };
    onMounted(async () => {
      const today = new Date();
      const nextThreeDays = new Date(today.setDate(today.getDate() + 3));
      const startDate = format(nextThreeDays, "yyyy-MM-dd");
      await fetchDateSlots();
      const events = dateSlots.value.map((slot) => {
        if (slot.available <= 0) {
          return {
            id: slot.id,
            title: "Fully Booked",
            start: slot.date,
            className: "p-2 bg-danger border border-none",
          };
        }
        return {
          id: slot.id,
          title: `Available Slots: ${slot.available ?? 0}`,
          start: slot.date,
          className: "p-2  bg-success cursor-pointer border border-none",
        };
      });

      reservationCalendar.value = new Calendar(
        reservationCalendarElement.value,
        {
          initialView: "dayGridMonth",
          plugins: [interactionPlugin],
          height: "650px",
          selectable: true,
          allDaySlot: false,
          validRange: {
            start: startDate,
          },
          eventClick: async (info) => {
            fetchTimeSlotsBasedOnDateSlot(info.event.id);
          },
          events: events,
        }
      );
      reservationCalendar.value.render();
    });
    return { reservationCalendarElement };
  },
}).mount("#ReservationPage");
