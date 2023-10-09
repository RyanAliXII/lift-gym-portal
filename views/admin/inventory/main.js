import { createApp } from "vue";
import { useForm } from "vee-validate";
import { format } from "date-fns";
import { object } from "yup";
import swal from "sweetalert2";
createApp({
  setup() {
    const { values, errors, defineInputBinds, setErrors } = useForm({
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
    const onSubmit = async () => {
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
        swal.fire("New equipment", "New equipment has been added.", "success");
      } catch (error) {
        console.error(error);
      }
    };
    return {
      name,
      model,
      quantity,
      costPrice,
      dateReceived,
      errors,
      onSubmit,
    };
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
}).mount("#InventoryPage");
