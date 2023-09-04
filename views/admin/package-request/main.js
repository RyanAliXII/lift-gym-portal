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
    const formatDate = (d) => {
      const date = new Date(d).toLocaleString({
        hourCycle: "h12",
      });
      return date;
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
    const initMarkAsReceived = async (id) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, mark it as received.",
        title: "Recieve Package Request",
        text: "Are you sure you want to mark request as received?",
        confirmButtonColor: "#295ad6",
        cancelButtonText: "I don't want to mark the request.",
        icon: "question",
      });
      if (result.isConfirmed) {
        updateStatus(id, Status.Received, async () => {
          swal.fire(
            "Package Request Receiving",
            "Package has been received by client.",
            "success"
          );
          fetchPackageRequests();
        });
      }
    };

    const initCancellation = async (id) => {
      const { value: text, isConfirmed } = await swal.fire({
        input: "textarea",
        title: "Cancel Package Request",
        inputLabel: "Cancellation Remarks",
        confirmButtonText: "Proceed to cancellation",
        cancelButtonText: "I don't want to cancel the request.",
        confirmButtonColor: "#d9534f",
        inputPlaceholder:
          "Enter cancellation reason eg. duplicate request, out of stock etc.",
        inputAttributes: {
          "aria-label": "Type your message here",
        },
        showCancelButton: true,
      });
      if (isConfirmed) {
        const formData = new FormData();
        formData.append("remarks", text);
        updateStatus(
          id,
          Status.Cancelled,
          async () => {
            swal.fire(
              "Package Request Cancellation",
              "Package request has been cancelled.",
              "success"
            );
            fetchPackageRequests();
          },
          formData
        );
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
      initMarkAsReceived,
      formatDate,
    };
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
}).mount("#PackageRequest");
