import { computed, createApp, onMounted, ref } from "vue";
import { Calendar } from "fullcalendar";
import interactionPlugin from "@fullcalendar/interaction";
import { format, parse, intervalToDuration } from "date-fns";
import swal from "sweetalert2";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const ReservationStatus = {
      Pending: 1,
      Attended: 2,
      NoShow: 3,
      Cancelled: 4,
    };
    const reservationCalendarElement = ref(null);
    const reservationCalendar = ref(null);
    const dateSlots = ref([]);
    const timeSlots = ref([]);
    const selectedDate = ref("");
    const reservations = ref([]);
    const initialValues = {
      dateSlotId: 0,
      timeSlotId: 0,
    };
    const errors = ref({});
    const form = ref({
      ...initialValues,
    });
    const reservationCache = computed(() => {
      const map = new Map();
      reservations.value.forEach((reservation) => {
        map.set(reservation.date, reservation);
      });
      return map;
    });
    const handleFormInput = (event) => {
      let value = event.target.value;
      let name = event.target.name;
      if (event.target.type === "number") {
        value = Number(value);
      }
      if (!isNaN(value)) {
        value = Number(value);
      }
      form.value[name] = value;
      delete errors.value[name];
    };
    const fetchDateSlots = async () => {
      try {
        const response = await fetch("/clients/reservations/date-slots");
        const { data } = await response.json();
        dateSlots.value = data?.slots ?? [];
      } catch (error) {
        console.error(error);
      }
    };
    const fetchReservations = async () => {
      try {
        const response = await fetch("/clients/reservations", {
          headers: new Headers({ "Content-Type": "application/json" }),
        });
        const { data } = await response.json();
        reservations.value = data?.reservations ?? [];
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
      await fetchReservations();
      await fetchDateSlots();

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
            if (reservationCache.value.has(slot.date)) return;
            selectedDate.value = formatDate(info.event.start);
            await fetchTimeSlotsBasedOnDateSlot(info.event.id);
            form.value.dateSlotId = parseInt(info.event.id);
            $("#reserveModal").modal("show");
          },
        }
      );
      repopulateEvents();
      $("#reserveModal").on("hidden.bs.modal", () => {
        form.value = { ...initialValues };
        errors.value = {};
      });
      reservationCalendar.value.render();
    });

    const formatDate = (date) => {
      if (!date) return "";
      return new Date(date).toLocaleDateString("en-US", {
        month: "long",
        day: "2-digit",
        year: "numeric",
      });
    };
    const repopulateEvents = () => {
      reservationCalendar.value.getEvents().forEach((event) => {
        event.remove();
      });
      dateSlots.value.forEach((slot) => {
        let event = {};
        if (reservationCache.value.has(slot.date)) {
          event = {
            id: slot.id,
            title: "Already Reserved",
            start: slot.date,
            className: "p-2 border border-none",
            extendedProps: slot,
          };
        } else if (slot.available <= 0) {
          event = {
            id: slot.id,
            title: "Fully Booked",
            start: slot.date,
            className: "p-2 bg-danger border border-none",
            extendedProps: slot,
          };
        } else {
          event = {
            id: slot.id,
            title: `Available Slots: ${slot.available ?? 0}`,
            start: slot.date,
            className: "p-2  bg-success cursor-pointer border border-none",
            extendedProps: slot,
          };
        }

        reservationCalendar.value.addEvent(event);
      });
    };
    const onSubmit = async () => {
      try {
        errors.value = {};
        const response = await fetch("/clients/reservations", {
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
        form.value = { ...initialValues };
        swal.fire(
          "New Reservation",
          "You have successfully reserved a slot.",
          "success"
        );
        await fetchDateSlots();
        await fetchReservations();
        repopulateEvents();
        $("#reserveModal").modal("hide");
      } catch (error) {
        console.error(error);
      }
    };

    const cancelReservation = async (id, remarks) => {
      const form = new FormData();
      if (remarks.length >= 0) {
        form.set("remarks", remarks);
      }
      const response = await fetch(
        `/clients/reservations/${id}/status/cancellation`,
        {
          method: "PUT",
          body: form,
          headers: new Headers({
            "X-CSRF-Token": window.csrf,
          }),
        }
      );
      if (response.status === 200) {
        swal.fire(
          "Reservation Status",
          "Reservation status has been cancelled.",
          "success"
        );

        fetchReservations();
      }
    };
    const isCancellable = (reservation) => {
      try {
        if (reservation.statusId != ReservationStatus.Pending) return false;
        let now = new Date();
        now.setHours(0, 0, 0, 0);
        let reservationDate = new Date(reservation.date);
        reservationDate.setHours(0, 0, 0, 0);
        const duration = intervalToDuration({
          start: reservationDate,
          end: now,
        });
        if (duration.days <= 0) {
          return false;
        }
        return true;
      } catch (error) {
        console.error(error);
        return false;
      }
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
      cancelReservation(id, text);
    };
    return {
      reservationCalendarElement,
      selectedDate,
      timeSlots,
      formatTime,
      onSubmit,
      form,
      handleFormInput,
      errors,
      reservations,
      isCancellable,
      formatDate,
      initCancellation,
    };
  },
}).mount("#ReservationPage");
