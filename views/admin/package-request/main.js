import { createApp, onMounted, ref } from "vue";
import swal from "sweetalert2";

const updateStatus = async (
  id,
  status = 0,
  onSuccess = () => {},
  data = new FormData()
) => {
  try {
    const response = await fetch(
      `/app/package-requests/${id}/status?statusId=${status}`,
      {
        method: "PATCH",
        body: data,
        headers: new Headers({
          "X-CSRF-Token": window.csrf,
        }),
      }
    );
    if (response.status === 200) {
      onSuccess();
    }
  } catch (error) {
    console.error(error);
  }
};

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

    const initApproval = async (id) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, approve it.",
        title: "Approve Package Request",
        text: "Are you sure you want to approve the package request?",
        confirmButtonColor: "#295ad6",
        cancelButtonText: "I don't want to approve the request.",
        icon: "question",
      });
      if (result.isConfirmed) {
        updateStatus(id, Status.Approved, async () => {
          swal.fire(
            "Package Request Approval",
            "Package request has been approved.",
            "success"
          );
          fetchPackageRequests();
        });
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
      initApproval,
    };
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
}).mount("#PackageRequest");
