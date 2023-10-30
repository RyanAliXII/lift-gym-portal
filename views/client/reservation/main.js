import { createApp, onMounted, ref } from "vue";
import { Calendar } from "fullcalendar";
import interactionPlugin from "@fullcalendar/interaction";
import { format, parse } from "date-fns";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const reservationCalendarElement = ref(null);
    const reservationCalendar = ref(null);
    const dateSlots = ref([]);
    const timeSlots = ref([]);
    const selectedDate = ref("");
    const initialValues = {
      dateSlotId: 0,
      timeSlotId: 0,
    };
    const form = ref({
      ...initialValues,
    });
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
            extendedProps: slot,
          };
        }
        return {
          id: slot.id,
          title: `Available Slots: ${slot.available ?? 0}`,
          start: slot.date,
          className: "p-2  bg-success cursor-pointer border border-none",
          extendedProps: slot,
        };
      });
      const formatDate = (date) => {
        if (!date) return "No Date";
        if (date.length === 0) return "No Date";
        return date.toLocaleDateString("en-US", {
          month: "long",
          day: "2-digit",
          year: "numeric",
        });
      };
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
            const slot = info.event.extendedProps;
            if (slot.available <= 0) return;
            selectedDate.value = formatDate(info.event.start);
            await fetchTimeSlotsBasedOnDateSlot(info.event.id);
            $("#reserveModal").modal("show");
          },
          events: events,
        }
      );
      reservationCalendar.value.render();
    });

    const onSubmit = () => {};
    return {
      reservationCalendarElement,
      selectedDate,
      timeSlots,
      formatTime,
      onSubmit,
      form,
    };
  },
}).mount("#ReservationPage");
