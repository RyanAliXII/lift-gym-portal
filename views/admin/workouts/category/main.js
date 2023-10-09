import { useForm } from "vee-validate";
import { createApp } from "vue";
import { object } from "yup";

createApp({
  setup() {
    const { values, defineInputBinds } = useForm({
      initialValues: {
        id: 0,
        name: "",
      },
      validationSchema: object({}),
    });
    const name = defineInputBinds("name", { validateOnChange: true });
    return {
      name,
    };
  },
}).mount("#WorkoutCategoryPage");
