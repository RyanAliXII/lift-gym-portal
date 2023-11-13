import { createApp, onMounted, ref } from "vue";
import swal from "sweetalert2";
import Choices from "choices.js";
import { useDebounceFn } from "@vueuse/core";
const updateStatus = async (
  id,
  status = 0,
  onSuccess = () => {},
  data = new FormData()
) => {
  try {
    const response = await fetch(
      `/app/package-requests/${id}/status?statusId=${status}`,
      {
        method: "PATCH",
        body: data,
        headers: new Headers({
          "X-CSRF-Token": window.csrf,
        }),
      }
    );
    if (response.status === 200) {
      onSuccess();
    }
  } catch (error) {
    console.error(error);
  }
};

createApp({
  setup() {
    const packageSelectElement = ref(null);
    const clientSelectElement = ref(null);
    let packageSelect = null;
    let clientSelect = null;
    const packageRequests = ref([]);
    const Status = {
      Pending: 1,
      Approved: 2,
      Received: 3,
      Cancelled: 4,
    };
    const form = ref({
      clientId: 0,
      packageId: 0,
    });
    const errors = ref({});
    const fetchPackageRequests = async () => {
      try {
        const response = await fetch("/app/package-requests", {
          headers: new Headers({ "content-type": "application/json" }),
        });
        const { data } = await response.json();
        packageRequests.value = data?.packageRequests ?? [];
      } catch (error) {
        console.error(error);
        packageRequests.value = [];
      }
    };
    const newPackageRequests = async () => {
      try {
        errors.value = {};
        const response = await fetch("/app/package-requests", {
          body: JSON.stringify(form.value),
          method: "POST",
          headers: new Headers({
            "Content-Type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
        });
        if (response.status === 400) {
          const { data } = await response.json();
          if (data?.errors) {
            errors.value = data?.errors ?? {};
          }
          return;
        }

        $("#requestModal").modal("hide");
        swal.fire(
          "Package Request Created",
          "Package request has been created.",
          "success"
        );
        fetchPackageRequests();
      } catch (error) {
        console.error(error);
      }
    };
    const fetchPackages = async () => {
      try {
        const response = await fetch("/app/packages", {
          headers: new Headers({
            "content-type": "application/json",
            "Cache-Control": "no-cache",
          }),
        });
        const { data } = await response.json();
        return data?.packages ?? [];
      } catch {
        return [];
      }
    };
    const formatDate = (d) => {
      const date = new Date(d).toLocaleString({
        hourCycle: "h12",
      });
      return date;
    };
    const initApproval = async (id) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, approve it.",
        title: "Approve Package Request",
        text: "Are you sure you want to approve the package request?",
        confirmButtonColor: "#295ad6",
        cancelButtonText: "I don't want to approve the request.",
        icon: "question",
      });
      if (result.isConfirmed) {
        updateStatus(id, Status.Approved, async () => {
          swal.fire(
            "Package Request Approval",
            "Package request has been approved.",
            "success"
          );
          fetchPackageRequests();
        });
      }
    };
    const initMarkAsReceived = async (id) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, mark it as received.",
        title: "Recieve Package Request",
        text: "Are you sure you want to mark request as received?",
        confirmButtonColor: "#295ad6",
        cancelButtonText: "I don't want to mark the request.",
        icon: "question",
      });
      if (result.isConfirmed) {
        updateStatus(id, Status.Received, async () => {
          swal.fire(
            "Package Request Receiving",
            "Package has been received by client.",
            "success"
          );
          fetchPackageRequests();
        });
      }
    };

    const initCancellation = async (id) => {
      const { value: text, isConfirmed } = await swal.fire({
        input: "textarea",
        title: "Cancel Package Request",
        inputLabel: "Cancellation Remarks",
        confirmButtonText: "Proceed to cancellation",
        cancelButtonText: "I don't want to cancel the request.",
        confirmButtonColor: "#d9534f",
        inputPlaceholder:
          "Enter cancellation reason eg. duplicate request, out of stock etc.",
        inputAttributes: {
          "aria-label": "Type your message here",
        },
        showCancelButton: true,
      });
      if (isConfirmed) {
        const formData = new FormData();
        formData.append("remarks", text);
        updateStatus(
          id,
          Status.Cancelled,
          async () => {
            swal.fire(
              "Package Request Cancellation",
              "Package request has been cancelled.",
              "success"
            );
            fetchPackageRequests();
          },
          formData
        );
      }
    };
    const fetchClientByKeyword = async (query) => {
      const response = await fetch(
        `/app/clients?${new URLSearchParams({
          keyword: query,
        }).toString()}`,
        {
          headers: new Headers({
            "Content-Type": "application/json",
            "Cache-Control": "no-cache",
          }),
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
        clientSelect.setChoices(selectValues, "value", "label", true);
      }
    };
    const search = useDebounceFn(fetchClientByKeyword, 500);
    const init = async () => {
      packageSelect = new Choices(packageSelectElement.value, {
        allowHTML: false,
      });
      clientSelect = new Choices(clientSelectElement.value, {
        allowHTML: false,
        placeholder: "Seach Client",
      });

      const packages = await fetchPackages();
      const packageOptions = packages.map((p) => ({
        value: p.id,
        label: p.description + " - " + p.price + " PHP",
        id: p.id,
        customProperties: p,
      }));

      clientSelect.passedElement.element.addEventListener("search", (event) => {
        search(event.detail.value);
      });
      clientSelect.passedElement.element.addEventListener("change", (event) => {
        delete errors.value.clientId;
        form.value.clientId = event.detail.value;
      });
      packageSelect.passedElement.element.addEventListener(
        "change",
        (event) => {
          delete errors.value.packageId;
          form.value.packageId = event.detail.value;
        }
      );
      packageSelect.setChoices(packageOptions);
      fetchPackageRequests();
      $("#requestModal").on("hidden.bs.modal", async () => {
        packageSelect.removeActiveItems();
        clientSelect.removeActiveItems();
      });
    };

    onMounted(() => {
      init();
    });

    return {
      packageRequests,
      Status,
      initCancellation,
      initApproval,
      initMarkAsReceived,
      formatDate,
      packageSelect,
      packageSelectElement,
      clientSelectElement,
      form,
      errors,
      newPackageRequests,
      errors,
    };
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
}).mount("#PackageRequest");
