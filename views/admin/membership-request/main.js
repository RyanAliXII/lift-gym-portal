import { createApp, ref, onMounted } from "vue";
import swal from "sweetalert2";

const fetchMembershipRequests = async () => {
  try {
    const response = await fetch("/app/membership-requests", {
      method: "GET",
      headers: new Headers({ "Content-Type": "application/json" }),
    });
    const { data } = await response.json();
    return data.membershipRequests ?? [];
  } catch (error) {
    return [];
    console.error(error);
  }
};

const updateStatus = async (
  id,
  status = 0,
  onSuccess = () => {},
  data = new FormData()
) => {
  try {
    const response = await fetch(
      `/app/membership-requests/${id}/status?statusId=${status}`,
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
    const membershipRequests = ref([]);
    const Status = {
      Pending: 1,
      Approved: 2,
      Received: 3,
      Cancelled: 4,
    };
    const initApproval = async (id) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, approve it.",
        title: "Approve Membership Request",
        text: "Are you sure you want to approve the membership request?",
        confirmButtonColor: "#295ad6",
        cancelButtonText: "I don't want to approve the request.",
        icon: "question",
      });
      if (result.isConfirmed) {
        updateStatus(id, Status.Approved, async () => {
          swal.fire(
            "Membership Request Approved",
            "Membership request has been approved.",
            "success"
          );
          membershipRequests.value = await fetchMembershipRequests();
        });
      }
    };
    const initCancellation = async (id) => {
      const { value: text, isConfirmed } = await swal.fire({
        input: "textarea",
        title: "Cancel Membership Request",
        inputLabel: "Cancellation Remarks",
        confirmButtonText: "Proceed to cancellation",
        cancelButtonText: "I don't want to cancel the request.",
        confirmButtonColor: "#d9534f",
        inputPlaceholder:
          "Enter cancellation reason eg. duplicate request etc.",
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
              "Membership Request Cancellation",
              "Membership request has been cancelled.",
              "success"
            );
            membershipRequests.value = await fetchMembershipRequests();
          },
          formData
        );
      }
    };
    const initMarkAsReceived = async (id) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, mark it as received.",
        title: "Recieve Membership Request",
        text: "Are you sure you want to mark request as received?",
        confirmButtonColor: "#295ad6",
        cancelButtonText: "I don't want to mark the request.",
        icon: "question",
      });
      if (result.isConfirmed) {
        updateStatus(id, Status.Received, async () => {
          swal.fire(
            "Membership Request Receiving",
            "Membership has been received by client.",
            "success"
          );
          membershipRequests.value = await fetchMembershipRequests();
        });
      }
    };

    const init = async () => {
      membershipRequests.value = await fetchMembershipRequests();
    };
    const formatDate = (d) => {
      const date = new Date(d).toLocaleString({
        hourCycle: "h12",
      });
      return date;
    };
    onMounted(() => {
      init();
    });
    return {
      membershipRequests,
      Status,
      initApproval,
      initMarkAsReceived,
      initCancellation,
      formatDate,
    };
  },
  compilerOptions: { delimiters: ["{", "}"] },
}).mount("#MembershipRequest");
