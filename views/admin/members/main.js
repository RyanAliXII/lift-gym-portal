import { createApp, onMounted } from "vue";
import Choices from "choices.js";
createApp({
  setup() {
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
      const selectPlan = new Choices("#selectPlan", {
        choices: planOptions,
      });
      const selectClient = new Choices("#selectClient", {
        choices: clientOptions,
      });
    };
    onMounted(() => {
      init();
    });
    return {};
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
}).mount("#MembersPage");
