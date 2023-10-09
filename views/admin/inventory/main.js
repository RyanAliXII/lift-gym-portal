import { createApp } from "vue";
import { useForm } from "vee-validate";
import { format } from "date-fns";
import { object } from "yup";
createApp({
  setup() {
    const {
      values: form,
      errors,
      defineInputBinds,
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
    const name = defineInputBinds(
      { value: "name" },
      { validateOnChange: true }
    );
    const model = defineInputBinds(
      { value: "model" },
      {
        validateOnChange: true,
      }
    );
    const quantity = defineInputBinds(
      { value: "quantity" },
      { validateOnChange: true }
    );
    const costPrice = defineInputBinds(
      { value: "costPrice" },
      { validateOnChange: true }
    );
    const dateReceived = defineInputBinds(
      { value: "dateReceived" },
      { validateOnChange: true }
    );
    const onSubmit = () => {};
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
