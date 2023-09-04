import { createApp, onMounted, ref } from "vue";
import Choices from "choices.js";
import { number } from "yup";
import swal from "sweetalert2";
const fetchMembershipPlans = async () => {
  try {
    const response = await fetch("/clients/memberships", {
      method: "GET",
      headers: new Headers({ "Content-Type": "application/json" }),
    });
    const { data } = await response.json();
    return data?.membershipPlans ?? [];
  } catch (error) {
    console.error(error);
    return [];
  }
};
const sendRequest = async (id = 0, onSuccess = () => {}) => {
  try {
    const response = await fetch("/clients/membership-requests", {
      method: "POST",
      body: JSON.stringify({
        membershipPlanId: id,
      }),
      headers: new Headers({
        "Content-Type": "application/json",
        "X-CSRF-Token": window.csrf,
      }),
    });
    if (response.status === 200) {
      onSuccess();
    }
  } catch (error) {
    console.error(error);
  }
};
const fetchMembershipRequests = async () => {
  try {
    const response = await fetch("/clients/membership-requests", {
      method: "GET",
      headers: new Headers({ "Content-Type": "application/json" }),
    });
    const { data } = await response.json();
    return data?.membershipRequests ?? [];
  } catch (error) {
    console.error(error);
    return [];
  }
};
createApp({
  setup() {
    const planSelectElement = ref(null);
    let planSelect = null;
    const errorMessage = ref(undefined);
    const membershipRequests = ref([]);
    const Status = {
      Pending: 1,
      Approved: 2,
      Received: 3,
      Cancelled: 4,
    };
    const init = async () => {
      const plans = await fetchMembershipPlans();
      const planOptions = plans.map((p) => ({
        value: p.id,
        label: p.description,
        id: p.id,
        customProperties: p,
      }));
      if (planSelectElement.value) {
        planSelect = new Choices(planSelectElement.value, {
          choices: planOptions,
        });
      }

      membershipRequests.value = await fetchMembershipRequests();
      $("#requestModal").on("shown.bs.modal", async () => {
        const plans = await fetchMembershipPlans();
        const planOptions = plans.map((p) => ({
          value: p.id,
          label: p.description,
          customProperties: p,
        }));
        planSelect.clearStore();
        planSelect.setChoices(planOptions);
      });
    };

    const onSubmit = async () => {
      errorMessage.value = undefined;
      const plan = planSelect.getValue();
      try {
        const result = await number()
          .required()
          .min(1)
          .validate(plan?.value ?? 0, { abortEarly: true });
        sendRequest(result, async () => {
          swal.fire(
            "Membership Request",
            "Membership request has been submitted. Please wait for response.",
            "success"
          );
          membershipRequests.value = await fetchMembershipRequests();
        });
      } catch (err) {
        errorMessage.value = "Please select a plan";
      } finally {
        $("#requestModal").modal("hide");
      }
    };
    const initCancellation = async (id) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, cancel it.",
        title: "Cancel Membership Request",
        text: "Are you sure you want to cancel the membership request?",
        confirmButtonColor: "#d9534f",
        cancelButtonText: "I don't want to cancel the request.",
        icon: "warning",
      });
      if (result.isConfirmed) {
        cancelMembershipRequest(id);
      }
    };
    const cancelMembershipRequest = async (id) => {
      try {
        const response = await fetch(
          `/clients/membership-requests/${id}/status?statusId=${Status.Cancelled}`,
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
            "Membership Request Cancellation",
            "Membership request has been cancelled.",
            "success"
          );
          membershipRequests.value = await fetchMembershipRequests();
        }
      } catch (error) {
        console.error(error);
      }
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
      planSelectElement,
      onSubmit,
      errorMessage,
      membershipRequests,
      Status,
      initCancellation,
      formatDate,
    };
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
}).mount("#MembershipRequest");
