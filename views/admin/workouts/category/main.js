import { useForm } from "vee-validate";
import { createApp } from "vue";
import { object } from "yup";
import swal from "sweetalert2";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const { values, defineInputBinds, errors, setErrors, resetForm } = useForm({
      initialValues: {
        id: 0,
        name: "",
      },
      validationSchema: object({}),
    });
    const onSubmitNew = async () => {
      try {
        const response = await fetch("/app/workouts/categories", {
          method: "POST",
          body: JSON.stringify(values),
          headers: new Headers({
            "Content-Type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
        });
        const { data } = await response.json();
        if (response.status >= 400) {
          if (data?.errors) {
            setErrors(data.errors);
          }
          return;
        }

        swal.fire(
          "Workout Category",
          "Workout category has been added.",
          "success"
        );
        $("#addCategoryModal").modal("hide");
        resetForm();
      } catch (error) {
        console.error(error);
      }
    };
    const name = defineInputBinds("name", { validateOnChange: true });
    return {
      name,
      errors,
      onSubmitNew,
    };
  },
}).mount("#WorkoutCategoryPage");
