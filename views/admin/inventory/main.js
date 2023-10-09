import { createApp } from "vue";
import { useForm } from "vee-validate";
import { format } from "date-fns";
import { object } from "yup";
createApp({
  setup() {
    const { values, errors, defineInputBinds, validate } = useForm({
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
      const form = await validate();
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
