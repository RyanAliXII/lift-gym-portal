import { createApp, onMounted, ref } from "vue";
import { useForm } from "vee-validate";
import { format } from "date-fns";
import { object } from "yup";
import swal from "sweetalert2";
createApp({
  setup() {
    const equipments = ref([]);
    const stat = ref({ totalCost: 0 });
    const {
      values,
      errors,
      defineInputBinds,
      setErrors,
      resetForm,
      setValues,
    } = useForm({
      initialValues: {
        name: "",
        model: "",
        quantity: 0,
        costPrice: 0,
        dateReceived: format(new Date(), "yyyy-MM-dd"),
      },
      validationSchema: object({}),
    });
    const name = defineInputBinds("name", { validateOnChange: true });
    const model = defineInputBinds("model", {
      validateOnChange: true,
    });
    const quantity = defineInputBinds("quantity", { validateOnChange: true });
    const costPrice = defineInputBinds("costPrice", { validateOnChange: true });
    const dateReceived = defineInputBinds("dateReceived", {
      validateOnChange: true,
    });
    const onSubmitNew = async () => {
      try {
        const response = await fetch("/app/inventory", {
          method: "POST",
          headers: new Headers({
            "Content-Type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
          body: JSON.stringify(values),
        });
        const { data } = await response.json();
        if (response.status === 400) {
          if (data?.errors) {
            setErrors(data?.errors);
          }
          return;
        }
        if (response.status === 200) {
          resetForm();
          $("#addEquipmentModal").modal("hide");
          fetchEquipments();
          swal.fire(
            "New equipment",
            "New equipment has been added.",
            "success"
          );
        }
      } catch (error) {
        console.error(error);
      }
    };

    const fetchEquipments = async () => {
      try {
        const response = await fetch("/app/inventory", {
          method: "GET",
          headers: new Headers({
            "Content-Type": "application/json",
            "X-CSRF-Token": window.csrf,
            "Cache-Control": "no-cache",
          }),
        });
        const { data } = await response.json();
        equipments.value = data?.equipments ?? [];
        stat.value = data?.stat ?? { totalCost: 0 };
      } catch (error) {
        console.error(error);
      }
    };

    const initEdit = (equipment) => {
      setValues(equipment);
      $("#editEquipmentModal").modal("show");
    };
    onMounted(() => {
      fetchEquipments();
      $("#addEquipmentModal").on("hidden.bs.modal", function () {
        resetForm();
      });
      $("#editEquipmentModal").on("hidden.bs.modal", function () {
        resetForm();
      });
    });
    const onSubmitUpdate = async () => {
      try {
        const url = `/app/inventory/${values?.id}`;
        const response = await fetch(url, {
          body: JSON.stringify(values),
          method: "PUT",
          headers: new Headers({
            "Content-Type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
        });
        const { data } = await response.json();
        if (response.status === 400) {
          if (data?.errors) {
            setErrors(data?.errors);
          }
          return;
        }
        if (response.status === 200) {
          resetForm();
          $("#editEquipmentModal").modal("hide");
          fetchEquipments();
          swal.fire(
            "Update equipment",
            "Equipment has been updated.",
            "success"
          );
        }
      } catch (error) {
        console.error(error);
      }
    };
    const deleteEquipment = async (id) => {
      const url = `/app/inventory/${id}`;
      const response = await fetch(url, {
        method: "DELETE",
        headers: new Headers({
          "Content-Type": "application/json",
          "X-CSRF-Token": window.csrf,
        }),
      });
      if (response.status === 200) {
        fetchEquipments();
        swal.fire("Delete equipment", "Equipment has been deleted.", "success");
      }
    };
    const initDelete = async (id) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, delete it",
        title: "Delete equipment",
        text: "Are you sure you want to delete equipment?",
        confirmButtonColor: "#d9534f",
        cancelButtonText: "I don't want to delete the equipment",
        icon: "warning",
      });
      if (result.isConfirmed) {
        deleteEquipment(id);
      }
    };
    const formatCurrency = (money) => {
      if (!money) return 0;
      return money.toLocaleString(undefined, {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
      });
    };
    return {
      name,
      model,
      quantity,
      costPrice,
      dateReceived,
      errors,
      equipments,
      onSubmitNew,
      onSubmitUpdate,
      initDelete,
      initEdit,
      stat,
      formatCurrency,
    };
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
}).mount("#InventoryPage");
