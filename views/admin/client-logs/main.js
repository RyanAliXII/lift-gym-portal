import Choices from "choices.js";
import { createApp, onMounted, ref } from "vue";
import { useDebounce, useDebounceFn } from "@vueuse/core";
import Swal from "sweetalert2";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const initialForm = {
      clientId: 0,
      isMember: false,
      amountPaid: 0,
    };
    const logs = ref([]);
    const logClientSelectElement = ref(null);
    const logClientSelect = ref(null);
    const form = ref({
      ...initialForm,
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
    const fetchLogs = async () => {
      try {
        const response = await fetch("/app/client-logs", {
          headers: new Headers({ "Content-Type": "application/json" }),
        });
        if (response.status >= 400) return;
        const { data } = await response.json();
        logs.value = data?.clientLogs ?? [];
      } catch (error) {
        console.error(error);
      }
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
      errors.value = {};
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
        if (response.status >= 400) {
          if (data?.errors) {
            errors.value = data?.errors;
          }
          return;
        }
        Swal.fire(
          "Client Loggged In",
          "Client has been loggged in successfully",
          "success"
        );
        form.value = {
          ...initialForm,
        };
        $("#logClientModal").modal("hide");
      } catch (error) {}
    };
    const initModalListeners = () => {
      $("#logClientModal").on("hidden.bs.modal", () => {
        logClientSelect.value.removeActiveItems();
        form.value = {
          ...initialForm,
        };
      });
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
      initModalListeners();
      fetchLogs();
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
