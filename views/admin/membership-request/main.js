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
    };
    const init = async () => {
      membershipRequests.value = await fetchMembershipRequests();
    };
    onMounted(() => {
      init();
    });
    return {
      membershipRequests,
      Status,
      initApproval,
    };
  },
  compilerOptions: { delimiters: ["{", "}"] },
}).mount("#MembershipRequest");
