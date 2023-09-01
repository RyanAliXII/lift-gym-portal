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
const sendRequest = async (id = 0) => {
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
      swal.fire(
        "Membership Request",
        "Membership request has been submitted. Please wait for response.",
        "success"
      );
    }
  } catch (error) {
    console.error(error);
  }
};
createApp({
  setup() {
    let planSelectElement = ref(null);
    let planSelect = null;
    let errorMessage = ref(undefined);

    const init = async () => {
      const plans = await fetchMembershipPlans();
      const planOptions = plans.map((p) => ({
        value: p.id,
        label: p.description,
        id: p.id,
        customProperties: p,
      }));
      planSelect = new Choices(planSelectElement.value, {
        choices: planOptions,
      });
    };

    const onSubmit = async () => {
      errorMessage.value = undefined;
      const id = planSelect.getValue()?.id;
      try {
        const result = await number()
          .required()
          .min(1)
          .validate(id, { abortEarly: true });
        sendRequest(result);
      } catch (err) {
        errorMessage.value = "Please select a plan";
      }
    };
    onMounted(() => {
      init();
    });
    return {
      planSelectElement,
      onSubmit,
      errorMessage,
    };
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
}).mount("#MembershipRequest");
