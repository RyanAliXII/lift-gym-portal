import { createApp, onMounted, ref } from "vue";
import Choices from "choices.js";
import { object, number } from "yup";
import { useForm } from "vee-validate";
import {
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
    let clientSelect = null;
    let planSelect = null;
    const {
      setValues,
      validate,
      errors,
      defineInputBinds,
      values: subscribeForm,
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
    const initCancellation = (member) => {
      swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, cancel it.",
        title: "Cancel Subscription",
        text: "Are you sure you want to cancel the subscription?",
        confirmButtonColor: "#d9534f",
        cancelButtonText: "I don't want to cancel the subscription",
        icon: "warning",
      });
    };
    const init = async () => {
      const plans = await fetchMembershipPlans();
      const clients = await fetchClients();
      members.value = await fetchMembers();
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
      planSelect = new Choices("#selectPlan", {
        choices: planOptions,
      });
      clientSelect = new Choices("#selectClient", {
        choices: clientOptions,
      });
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
    };
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
}).mount("#MembersPage");
