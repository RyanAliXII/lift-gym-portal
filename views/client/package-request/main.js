import { createApp, onMounted, ref } from "vue";
import Choices from "choices.js";
import { number } from "yup";
import swal from "sweetalert2";
createApp({
  setup() {
    const packageSelectElement = ref(null);
    let packageSelect = null;
    const packageRequests = ref([]);
    const errorMessage = ref(undefined);
    const Status = {
      Pending: 1,
      Approved: 2,
      Received: 3,
      Cancelled: 4,
    };
    const fetchPackages = async () => {
      try {
        const response = await fetch("/clients/packages", {
          headers: new Headers({ "content-type": "application/json" }),
        });
        const { data } = await response.json();
        return data?.packages ?? [];
      } catch {
        return [];
      }
    };

    const fetchPackageRequests = async () => {
      try {
        const response = await fetch("/clients/package-requests", {
          headers: new Headers({ "content-type": "application/json" }),
        });
        const { data } = await response.json();
        packageRequests.value = data?.packageRequests ?? [];
      } catch (error) {
        console.error(error);
        packageRequests.value = [];
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
          fetchPackageRequests();
          $("#requestModal").modal("hide");
        }
      } catch (error) {
        console.error(error);
      }
    };
    const initCancellation = async (id) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, cancel it.",
        title: "Cancel Package Request",
        text: "Are you sure you want to cancel the package request?",
        confirmButtonColor: "#d9534f",
        cancelButtonText: "I don't want to cancel the request.",
        icon: "warning",
      });
      if (result.isConfirmed) {
        cancelPackageRequest(id);
      }
    };

    const cancelPackageRequest = async (id) => {
      try {
        const response = await fetch(
          `/clients/package-requests/${id}/status?statusId=${Status.Cancelled}`,
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
            "Package Request Cancellation",
            "Package request has been cancelled.",
            "success"
          );
          fetchPackageRequests();
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
      fetchPackageRequests();
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
      packageRequests,
      Status,
      initCancellation,
    };
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
}).mount("#PackageRequest");
