import { createApp, onMounted, ref } from "vue";
import Choices from "choices.js";
import { number } from "yup";
import swal from "sweetalert2";
createApp({
  setup() {
    const packageSelectElement = ref(null);
    let packageSelect = null;

    const errorMessage = ref(undefined);

    const fetchPackages = async () => {
      try {
        const response = await fetch("/clients/packages");
        const { data } = await response.json();
        return data?.packages ?? [];
      } catch {
        return [];
      }
    };
    const onSubmit = async () => {
      try {
        const pkg = packageSelect.getValue();
        const result = await number().required().min(1).validate(pkg.value);
        sendRequest({ packageId: result });
      } catch (error) {
        errorMessage.value = "Package is required.";
        console.error(error);
      }
    };
    const sendRequest = async (
      form = {
        packageId: 0,
      }
    ) => {
      errorMessage.value = undefined;
      try {
        const response = await fetch("/clients/package-requests", {
          body: JSON.stringify(form),
          method: "POST",
          headers: new Headers({
            "Content-Type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
        });

        if (response.status === 200) {
          swal.fire(
            "Package Request",
            "Package request has been submitted. Please wait for response.",
            "success"
          );
          $("#requestModal").modal("hide");
        }
      } catch (error) {
        console.error(error);
      }
    };
    const init = async () => {
      packageSelect = new Choices(packageSelectElement.value, {});
      const packages = await fetchPackages();
      const packageOptions = packages.map((p) => ({
        value: p.id,
        label: p.description + " - " + p.price + " PHP",
        id: p.id,
        customProperties: p,
      }));
      packageSelect.setChoices(packageOptions);
    };
    onMounted(() => {
      init();

      $("#requestModal").on("shown.bs.modal", async () => {
        packageSelect.clearStore();
        const packages = await fetchPackages();
        const packageOptions = packages.map((p) => ({
          value: p.id,
          label: p.description + " - " + p.price + " PHP",
          id: p.id,
          customProperties: p,
        }));
        packageSelect.setChoices(packageOptions);
      });
    });

    return {
      packageSelectElement,
      errorMessage,
      onSubmit,
    };
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
}).mount("#PackageRequest");
