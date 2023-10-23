import { createApp, onMounted, ref } from "vue";
import swal from "sweetalert2";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const hiredCoaches = ref([]);
    const fetchHiredCoaches = async () => {
      try {
        const response = await fetch("/clients/hired-coaches", {
          headers: new Headers({ "Content-Type": "application/json" }),
        });
        const { data } = await response.json();

        if (response.status === 200) {
          hiredCoaches.value = data?.hiredCoaches ?? [];
        }
      } catch (error) {
        console.error(error);
      }
    };
    const cancelRequest = async (id) => {
      try {
        const response = await fetch(`/clients/hired-coaches/${id}`, {
          method: "DELETE",
          headers: new Headers({
            "Content-Type": "application/json",
            "X-CSRF-token": window.csrf,
          }),
        });
        if (response.status >= 400) return;

        swal.fire(
          "Coach Appoinment cancellation.",
          "Coach appointment has been cancelled.",
          "success"
        );
      } catch (err) {
        console.error(err);
      }
    };
    const initCancel = async (id) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, cancel it.",
        title: "Cancel Appointment",
        text: "Are you sure you want to cancel coaching request?",
        confirmButtonColor: "#d9534f",
        cancelButtonText: "I don't want to cancel coaching request",
        icon: "warning",
      });
      if (result.isDenied || result.isDismissed) return;
      cancelRequest(id);
    };
    const toMoney = (money) => {
      if (!money) return 0;
      return money.toLocaleString(undefined, {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
      });
    };
    onMounted(() => {
      fetchHiredCoaches();
    });
    return {
      hiredCoaches,
      toMoney,
      initCancel,
    };
  },
}).mount("#HiredCoaches");
