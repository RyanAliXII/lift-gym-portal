import { createApp, onMounted, ref } from "vue";
import { Calendar } from "fullcalendar";

import interactionPlugin from "@fullcalendar/interaction";
import { format } from "date-fns";
createApp({
  setup() {
    const reservationCalendarElement = ref(null);
    const reservationCalendar = ref(null);
    const dateSlots = ref([]);

    const fetchDateSlots = async () => {
      try {
        const response = await fetch("/clients/reservations/date-slots");
        const { data } = await response.json();
        dateSlots.value = data?.slots ?? [];
      } catch (error) {
        console.error(error);
      }
    };
    onMounted(async () => {
      const today = new Date();
      const nextThreeDays = new Date(today.setDate(today.getDate() + 3));
      const startDate = format(nextThreeDays, "yyyy-MM-dd");
      await fetchDateSlots();
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
          dateClick: (info) => {
            const date = info.dateStr;
            // if (allowedDates.includes(date)) {
            //   reservationCalendar.value.changeView("timeGridDay");
            // }
          },
          eventClick: (info) => {
            // reservationCalendar.value.changeView(
            //   "timeGridDay",
            //   info.event.startStr
            // );
          },
          events: dateSlots.value.map((slot) => ({
            title: "This date is open for reservation.",
            start: slot.date,
            className: "p-2  bg-success cursor-pointer",
          })),
        }
      );
      reservationCalendar.value.render();
    });
    return { reservationCalendarElement };
  },
}).mount("#ReservationPage");
