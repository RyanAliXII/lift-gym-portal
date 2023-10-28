import { format } from "date-fns";
import { createApp, onMounted, ref } from "vue";
import swal from "sweetalert2";
import { Calendar } from "fullcalendar";
import timeGridPlugin from "@fullcalendar/timegrid";
import interactionPlugin from "@fullcalendar/interaction";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const today = format(new Date(), "yyyy-MM-dd");
    const calendarViewElement = ref(null);
    const calendarView = ref(null);
    const dateWithSlots = ref([]);
    const initialFormValue = {
      to: today,
      from: today,
    };
    const form = ref({
      ...initialFormValue,
    });
    const slots = ref([]);
    const errors = ref({});
    const handleFormInput = (event) => {
      let value = event.target.value;
      let name = event.target.name;
      if (event.target.type === "number") {
        value = Number(value);
      }
      form.value[name] = value;
      delete errors.value[name];
    };
    const onSubmit = async () => {
      try {
        errors.value = {};
        const response = await fetch("/app/date-slots", {
          method: "POST",
          body: JSON.stringify(form.value),
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
        form.value = { ...initialFormValue };
        $("#newSlotModal").modal("hide");
        fetchSlots();
        swal.fire("New Slot", "Date slot/s has been added.", "success");
      } catch (error) {
        console.error(error);
      }
    };
    const fetchSlots = async () => {
      try {
        const response = await fetch("/app/date-slots", {
          headers: new Headers({ "Content-Type": "application/json" }),
        });
        const { data } = await response.json();
        slots.value = data?.slots ?? [];

        populateEvents(data?.slots ?? []);
      } catch (error) {
        console.error(error);
      }
    };

    const formatDate = (date) => {
      if (!date) return "No Date";
      if (date.length === 0) return "No Date";
      return new Date(date).toLocaleDateString("en-US", {
        month: "long",
        day: "2-digit",
        year: "numeric",
      });
    };
    const deleteSlot = async (id) => {
      try {
        const response = await fetch(`/app/date-slots/${id}`, {
          method: "DELETE",
          headers: new Headers({ "X-CSRF-Token": window.csrf }),
        });
        if (response.status === 200) {
          swal.fire(
            "Delete Date Slot",
            "Date slot has been deleted.",
            "success"
          );
          fetchSlots();
        }
      } catch (error) {
        console.error(error);
      }
    };
    const initDelete = async (id) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, delete it",
        title: "Delete Date Slot",
        text: "Are you sure you want to delete slot",
        confirmButtonColor: "#d9534f",
        cancelButtonText: "I don't want to delete this slot",
        icon: "warning",
      });
      if (result.isConfirmed) {
        deleteSlot(id);
      }
    };
    const populateEvents = (slots = []) => {
      clearEvents();
      dateWithSlots.value = [];
      slots.forEach((slot) => {
        dateWithSlots.value.push(slot.date);
        calendarView.value.addEvent({
          id: slot.id,
          title: "This date is available for reservation",
          date: slot.date,
          className: "p-2  bg-success",
        });
      });
    };

    onMounted(() => {
      const today = new Date();
      const startDate = format(today, "yyyy-MM-dd");

      calendarView.value = new Calendar(calendarViewElement.value, {
        initialView: "dayGridMonth",
        plugins: [timeGridPlugin, interactionPlugin],
        height: "650px",
        selectable: true,
        allDaySlot: false,

        headerToolbar: {
          left: "prev,next",
          center: "title",
          right: "dayGridMonth",
        },
        validRange: {
          start: startDate,
        },
        dateClick: (info) => {
          const date = info.dateStr;
          if (dateWithSlots.value.includes(date)) {
            calendarView.value.changeView("timeGridDay");
          }
        },
        eventClick: (info) => {
          calendarView.value.changeView("timeGridDay", info.event.startStr);
        },
      });

      $('a[data-toggle="tab"]').on("shown.bs.tab", (event) => {
        if (event.target.id === "calendar-view-tab") {
          calendarView.value.render();
          calendarView.value.changeView("dayGridMonth");
        }
      });
      fetchSlots();
    });

    const clearEvents = () => {
      calendarView.value.getEvents().forEach((event) => {
        event.remove();
      });
    };
    return {
      form,
      handleFormInput,
      onSubmit,
      errors,
      today,
      formatDate,
      slots,
      initDelete,
      calendarViewElement,
    };
  },
}).mount("#DateSlot");
