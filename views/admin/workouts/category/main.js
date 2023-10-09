import { useForm } from "vee-validate";
import { createApp, onMounted, ref } from "vue";
import { object } from "yup";
import swal from "sweetalert2";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const {
      values,
      defineInputBinds,
      errors,
      setErrors,
      resetForm,
      setValues,
    } = useForm({
      initialValues: {
        id: 0,
        name: "",
      },
      validationSchema: object({}),
    });
    const categories = ref([]);
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
        fetchCategories();
        $("#addCategoryModal").modal("hide");
        resetForm();
      } catch (error) {
        console.error(error);
      }
    };

    const fetchCategories = async () => {
      try {
        const response = await fetch("/app/workouts/categories", {
          headers: new Headers({ "Content-Type": "application/json" }),
        });
        if (response.status === 200) {
          const { data } = await response.json();
          categories.value = data?.categories ?? [];
        }
      } catch (error) {
        console.error(error);
      }
    };

    const onSubmitUpdate = async () => {
      try {
        const response = await fetch(`/app/workouts/categories/${values.id}`, {
          method: "PUT",
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
          "Workout category has been updated.",
          "success"
        );
        fetchCategories();
        $("#editCategoryModal").modal("hide");
      } catch (error) {
        console.error(error);
      }
    };
    const initEdit = (category) => {
      setValues(category);
      $("#editCategoryModal").modal("show");
    };
    onMounted(() => {
      fetchCategories();
      $("#addCategoryModal").on("hidden.bs.modal", function () {
        resetForm();
      });
      $("#editCategoryModal").on("hidden.bs.modal", function () {
        resetForm();
      });
    });
    const name = defineInputBinds("name", { validateOnChange: true });
    return {
      name,
      errors,
      categories,
      onSubmitNew,
      initEdit,
      onSubmitUpdate,
    };
  },
}).mount("#WorkoutCategoryPage");
