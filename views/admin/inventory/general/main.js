import { createApp, onMounted, ref } from "vue";
import swal from "sweetalert2";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const INITIAL_FORM = {
      name: "",
      brand: "",
      quantity: 0,
      unitOfMeasure: "",
      dateReceived: "",
      quantityThreshold: 0,
      costPrice: 0,
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

    onMounted(() => {
      fetchItems();
    });
    return {
      errors,
      form,
      handleImage,
      onSubmit,
      isLoading,
      formatCurrency,
      items,
    };
  },
}).mount("#GeneralInventoryPage");
