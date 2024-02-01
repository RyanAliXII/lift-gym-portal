import { createApp, onMounted, ref } from "vue";
import swal from "sweetalert2";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const INITIAL_FORM = {
      id: 0,
      name: "",
      brand: "",
      quantity: 0,
      unitOfMeasure: "",
      dateReceived: "",
      quantityThreshold: 0,
      costPrice: 0,
      image: "",
    };
    const isLoading = ref(false);
    const errors = ref({});
    const image = ref(null);
    const form = ref({ ...INITIAL_FORM });
    const items = ref([]);
    const fetchItems = async () => {
      try {
        const response = await fetch("/app/general/inventory", {
          headers: new Headers({
            "Content-Type": "application/json",
            "Cache-Control": "no-cache",
          }),
        });
        const { data } = await response.json();
        items.value = data?.items;
      } catch (error) {
        console.error(error);
      }
    };

    const onSubmit = async () => {
      try {
        isLoading.value = true;
        errors.value = {};
        const formData = toFormData();
        const response = await fetch("/app/general/inventory", {
          method: "POST",
          body: formData,
          headers: new Headers({
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
        form.value = { ...INITIAL_FORM };
        swal.fire(
          "Inventory Item",
          "New item has been added to invetory.",
          "success"
        );
        $("#addItemModal").modal("hide");
        fetchItems();
      } catch (error) {
      } finally {
        isLoading.value = false;
      }
    };

    const onSubmitUpdate = async () => {
      try {
        isLoading.value = true;
        errors.value = {};
        const formData = toFormData();
        const response = await fetch(
          `/app/general/inventory/${form.value.id}`,
          {
            method: "PUT",
            body: formData,
            headers: new Headers({
              "X-CSRF-Token": window.csrf,
            }),
          }
        );
        const { data } = await response.json();

        if (response.status >= 400) {
          if (data?.errors) {
            errors.value = data?.errors;
          }
          return;
        }
        form.value = { ...INITIAL_FORM };
        swal.fire("Inventory Item", "Item has been updated.", "success");
        $("#editItemModal").modal("hide");
        fetchItems();
      } catch (error) {
      } finally {
        isLoading.value = false;
      }
    };
    const toFormData = () => {
      const formData = new FormData();
      for (const [key, value] of Object.entries(form.value)) {
        formData.append(key, value);
      }
      if (image.value) {
        formData.append("imageFile", image.value);
      }
      return formData;
    };
    const handleImage = (event) => {
      const img = event.target.files[0];
      if (!img) return;
      image.value = img;
    };
    const formatCurrency = (money) => {
      if (!money) return 0;
      return money.toLocaleString(undefined, {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
      });
    };

    const initEdit = (item) => {
      errors.value = {};
      form.value = { ...item };
      $("#editItemModal").modal("show");
    };
    const initDelete = async (itemId) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, delete it",
        title: "Delete Inventory Item",
        text: "Are you sure you want to delete inventory item?",
        confirmButtonColor: "#d9534f",
        cancelButtonText: "I don't want to delete the item",
        icon: "warning",
      });
      if (result.isConfirmed) {
        deleteItem(itemId);
      }
    };
    const deleteItem = async (id) => {
      try {
        const response = await fetch(`/app/general/inventory/${id}`, {
          method: "DELETE",
          headers: new Headers({
            "Content-Type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
        });
        if (response.status === 200) {
          swal.fire(
            "Inventory Item",
            "Item deleted from inventory.",
            "success"
          );
          fetchItems();
        }
      } catch (error) {
        console.error(error);
      }
    };

    onMounted(() => {
      fetchItems();
    });
    return {
      errors,
      form,
      handleImage,
      onSubmit,
      isLoading,
      initEdit,
      formatCurrency,
      items,
      initDelete,
      onSubmitUpdate,
    };
  },
}).mount("#GeneralInventoryPage");
