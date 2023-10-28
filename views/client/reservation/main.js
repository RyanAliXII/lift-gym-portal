import { createApp, onMounted, ref } from "vue";
import { Calendar } from "fullcalendar";
import timeGridPlugin from "@fullcalendar/timegrid";
import interactionPlugin from "@fullcalendar/interaction";
import { format } from "date-fns";
import { da } from "date-fns/locale";
createApp({
  setup() {
    const reservationCalendarElement = ref(null);
    const reservationCalendar = ref(null);
    onMounted(() => {
      const today = new Date();
      const nextThreeDays = new Date(today.setDate(today.getDate() + 3));
      const startDate = format(nextThreeDays, "yyyy-MM-dd");
      const allowedDates = ["2023-11-05", "2023-11-02"];
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
            center: "title",
            right: "dayGridMonth", // user can switch between the two
          },
          validRange: {
            start: startDate,
          },
          dateClick: (info) => {
            const date = info.dateStr;
            if (allowedDates.includes(date)) {
              reservationCalendar.value.changeView("timeGridDay");
            }
          },
          eventClick: (info) => {
            reservationCalendar.value.changeView(
              "timeGridDay",
              info.event.startStr
            );
          },
          events: allowedDates.map((d) => ({
            title: "This date is open for reservation.",
            start: d,
            className: "p-2  bg-success",
          })),
        }
      );
      reservationCalendar.value.render();
    });
    return { reservationCalendarElement };
  },
}).mount("#ReservationPage");
