import Choices from "choices.js";
import { createApp, onMounted, ref } from "vue";
import { useDebounce, useDebounceFn } from "@vueuse/core";
createApp({
  setup() {
    const logClientSelectElement = ref(null);
    const logClientSelect = ref(null);
    const form = ref({
      clientId: 0,
      isMember: false,
      amountPaid: 0,
    });
    const search = useDebounceFn(async (query) => {
      const response = await fetch(
        `/app/clients?${new URLSearchParams({
          keyword: query,
        }).toString()}`,
        {
          headers: new Headers({ "Content-Type": "application/json" }),
        }
      );

      if (response.status === 200) {
        const { data } = await response.json();
        const selectValues = (data?.clients ?? []).map((client) => ({
          value: client.id,
          label: `${client.givenName} ${client.surname} - ${client.email} - ${
            client.isMember ? "Member" : "Non-Member"
          }`,
          customProperties: client,
        }));
        logClientSelect.value.setChoices(selectValues, "value", "label", true);
      }
    }, 500);

    onMounted(() => {
      logClientSelect.value = new Choices(logClientSelectElement.value, {
        allowHTML: false,
        placeholder: "Seach Client",
      });

      logClientSelect.value.passedElement.element.addEventListener(
        "search",
        (event) => {
          search(event.detail.value);
        }
      );
      logClientSelect.value.passedElement.element.addEventListener(
        "change",
        () => {
          const select = logClientSelect.value.getValue();

          form.value = {
            ...form.value,
            clientId: select.value,
            isMember: select.customProperties.isMember,
          };
        }
      );
    });

    return {
      logClientSelectElement,
      form,
    };
  },
}).mount("#ClientLog");
