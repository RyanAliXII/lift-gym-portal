import { createApp, onMounted } from "vue";
import { useForm } from "vee-validate";
import { object, string } from "yup";
createApp({
  setup() {
    const { values, defineInputBinds, errors, handleSubmit } = useForm({
      validationSchema: object({
        name: string().required(),
      }),
      initialValues: {
        name: "",
      },
    });

    onMounted(() => {
      console.log("App Mounted");
    });
    const name = defineInputBinds("name");
    const onSubmit = handleSubmit((values) => {
      console.log(values);
    });
    return {
      name,
      errors,
      onSubmit,
    };
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
}).mount("#MembersPage");
