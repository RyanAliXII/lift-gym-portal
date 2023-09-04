import { createApp, onMounted, ref } from "vue";
import Choices from "choices.js";
import { number } from "yup";
import swal from "sweetalert2";
createApp({
  setup() {
    const packageRequests = ref([]);
    const Status = {
      Pending: 1,
      Approved: 2,
      Received: 3,
      Cancelled: 4,
    };
    const fetchPackageRequests = async () => {
      try {
        const response = await fetch("/app/package-requests", {
          headers: new Headers({ "content-type": "application/json" }),
        });
        const { data } = await response.json();
        packageRequests.value = data?.packageRequests ?? [];
      } catch (error) {
        console.error(error);
        packageRequests.value = [];
      }
    };

    const initCancellation = async (id) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, cancel it.",
        title: "Cancel Package Request",
        text: "Are you sure you want to cancel the package request?",
        confirmButtonColor: "#d9534f",
        cancelButtonText: "I don't want to cancel the request.",
        icon: "warning",
      });
      if (result.isConfirmed) {
        cancelPackageRequest(id);
      }
    };

    const cancelPackageRequest = async (id) => {
      try {
        const response = await fetch(
          `/clients/package-requests/${id}/status?statusId=${Status.Cancelled}`,
          {
            method: "PATCH",
            headers: new Headers({
              "Content-Type": "application/json",
              "X-CSRF-Token": window.csrf,
            }),
          }
        );
        if (response.status === 200) {
          swal.fire(
            "Package Request Cancellation",
            "Package request has been cancelled.",
            "success"
          );
          fetchPackageRequests();
        }
      } catch (error) {
        console.error(error);
      }
    };
    const init = async () => {
      fetchPackageRequests();
    };
    onMounted(() => {
      init();
    });

    return {
      packageRequests,
      Status,
      initCancellation,
    };
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
}).mount("#PackageRequest");
