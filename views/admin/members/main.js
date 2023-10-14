import { createApp, onMounted, ref } from "vue";
import Choices from "choices.js";
import { object, number } from "yup";
import { useForm } from "vee-validate";
import {
  cancelSubscription,
  fetchClients,
  fetchMembers,
  fetchMembershipPlans,
  subscribe,
} from "./fetch";
import swal from "sweetalert2";

const SubscribeValidation = object({
  clientId: number()
    .required("Client is required")
    .integer("Client is required.")
    .min(1, "Client is required"),
  membershipPlanId: number()
    .required("Membership plan is required")
    .integer("Membership plan is required")
    .min(1, "Membership plan is required"),
});
createApp({
  setup() {
    const clientSelectElement = ref(null);
    const planSelectElement = ref(null);
    let clientSelect = null;
    let planSelect = null;
    const {
      setValues,
      validate,
      errors,
      defineInputBinds,
      values: subscribeForm,
      setFieldError,
      setErrors,
    } = useForm({
      initialValues: {
        clientId: 0,
        membershipPlanId: 0,
      },
      validateOnMount: false,
      validationSchema: SubscribeValidation,
    });

    const members = ref([]);
    const onSubmit = async () => {
      let client = clientSelect.getValue();
      let plan = planSelect.getValue();
      setValues({ clientId: client?.value, membershipPlanId: plan?.value });
      const { valid } = await validate();
      if (valid) {
        subscribe(subscribeForm, async () => {
          swal.fire(
            "Subscribe Client",
            `Client has been subscribed successfully.`,
            "success"
          );
          members.value = await fetchMembers();
        });
      }
    };
    const formatDate = (dateStr) => {
      return new Date(dateStr).toLocaleDateString("en-US", {
        month: "long",
        day: "2-digit",
        year: "numeric",
      });
    };

    const initCancellation = async (member) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, cancel it.",
        title: "Cancel Subscription",
        text: "Are you sure you want to cancel the subscription?",
        confirmButtonColor: "#d9534f",
        cancelButtonText: "I don't want to cancel the subscription",
        icon: "warning",
      });
      if (result.isConfirmed) {
        cancelSubscription(member.subscriptionId, async () => {
          swal.fire(
            "Cancel Subscription",
            "Subscription has been cancelled",
            "success"
          );
          members.value = await fetchMembers();
        });
      }
    };

    const initListeners = () => {
      $("#subscribeClientModal").on("shown.bs.modal", async () => {
        planSelect.clearStore();
        clientSelect.clearStore();
        setErrors({
          clientId: undefined,
          membershipPlanId: undefined,
        });
        const plans = await fetchMembershipPlans();
        const clients = await fetchClients();
        const planOptions = plans.map((p) => ({
          value: p.id,
          label: p.description,
          id: p.id,
          customProperties: p,
        }));
        const clientOptions = clients.map((c) => ({
          value: c.id,
          label: `${c.givenName} ${c.surname}`,
          customProperties: c,
        }));
        planSelect.setChoices(planOptions);
        clientSelect.setChoices(clientOptions);
      });

      planSelectElement.value.addEventListener("change", () => {
        if (errors.value.membershipPlanId) {
          setFieldError("membershipPlanId", undefined);
        }
      });
      clientSelectElement.value.addEventListener("change", () => {
        if (errors.value.clientId) {
          setFieldError("clientId", undefined);
        }
      });
    };
    const init = async () => {
      members.value = await fetchMembers();
      planSelect = new Choices(planSelectElement.value, {});
      clientSelect = new Choices(clientSelectElement.value, {});
      initListeners();
    };
    defineInputBinds("clientId");
    defineInputBinds("membershipPlanId");
    onMounted(() => {
      init();
    });
    return {
      onSubmit,
      errors,
      members,
      formatDate,
      initCancellation,
      clientSelectElement,
      planSelectElement,
    };
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
}).mount("#MembersPage");
