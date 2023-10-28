import { createApp, onMounted, ref } from "vue";
import { Calendar } from "fullcalendar";
import timeGridPlugin from "@fullcalendar/timegrid";
import interactionPlugin from "@fullcalendar/interaction";
createApp({
  setup() {
    const reservationCalendarElement = ref(null);
    const reservationCalendar = ref(null);
    onMounted(() => {
      reservationCalendar.value = new Calendar(
        reservationCalendarElement.value,
        {
          initialView: "dayGridMonth",
          plugins: [timeGridPlugin, interactionPlugin],
          height: "650px",
          selectable: true,
          allDaySlot: false,
          headerToolbar: {
            left: "prev,next",
            right: "dayGridMonth,timeGridDay", // user can switch between the two
          },
          dateClick: (info) => {
            reservationCalendar.value.changeView("timeGridDay", info.dateStr);
          },
        }
      );
      reservationCalendar.value.render();
    });
    return { reservationCalendarElement };
  },
}).mount("#ReservationPage");
