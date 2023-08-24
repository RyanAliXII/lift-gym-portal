import { createApp, onMounted } from "vue";
import Choices from "choices.js";
import { object, number } from "yup";
import { useForm } from "vee-validate";
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

    const fetchMembershipPlans = async () => {
      try {
        const response = await fetch("/memberships", {
          headers: new Headers({ "content-type": "application/json" }),
        });
        const { data } = await response.json();
        return data?.membershipPlans;
      } catch (error) {
        console.error(error);
        return [];
      }
    };
    const fetchClients = async () => {
      try {
        const response = await fetch("/clients", {
          headers: new Headers({ "content-type": "application/json" }),
        });
        const { data } = await response.json();
        return data?.clients;
      } catch (error) {
        console.error(error);
        return [];
      }
    };
    const onSubmit = async () => {
      let client = clientSelect.getValue();
      let plan = planSelect.getValue();
      setValues({ clientId: client?.value, membershipPlanId: plan?.value });
      console.log();
      const { valid } = await validate();
      if (valid) {
        subscribe();
      }
    };

    const subscribe = async () => {
      try {
        const response = await fetch("/members", {
          method: "POST",
          body: JSON.stringify(subscribeForm),
          headers: new Headers({
            "content-type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
        });
        if (response.status === 200) {
          swal.fire(
            "Subscribe Client",
            `Client has been subscribed successfully.`,
            "success"
          );
        }
      } catch (error) {}
    };
    const init = async () => {
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
    };
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
}).mount("#MembersPage");
