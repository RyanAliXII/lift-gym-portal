import Choices from "choices.js";
import { createApp, onMounted, ref } from "vue";
import { useDebounce, useDebounceFn } from "@vueuse/core";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const logClientSelectElement = ref(null);
    const logClientSelect = ref(null);
    const form = ref({
      clientId: 0,
      isMember: false,
      amountPaid: 0,
    });
    const errors = ref({});
    const handleFormInput = (event) => {
      let value = event.target.value;
      let name = event.target.name;
      if (event.target.type === "number") {
        value = Number(value);
      }
      form.value[name] = value;
      delete errors.value[name];
    };
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
    const submitLog = async () => {
      try {
        const response = await fetch("/app/client-logs", {
          method: "POST",
          body: JSON.stringify(form.value),
          headers: new Headers({
            "Content-Type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
        });
        const { data } = await response.json();
        if (response.status === 400) {
          if (data?.errors) {
            errors.value = data?.errors;
          }
          return;
        }
      } catch (error) {}
    };
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
          delete errors.clientId;
        }
      );
    });

    return {
      logClientSelectElement,
      form,
      handleFormInput,
      submitLog,
      errors,
    };
  },
}).mount("#ClientLog");