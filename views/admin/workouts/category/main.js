import { useForm } from "vee-validate";
import { createApp, onMounted, ref, watch } from "vue";
import { object } from "yup";
import swal from "sweetalert2";
import Choices from "choices.js";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const {
      values,
      defineInputBinds,
      errors,
      resetForm,
      setValues,
      setErrors,
    } = useForm({
      initialValues: {
        id: 0,
        name: "",
        workouts: [],
      },
      validationSchema: object({}),
    });

    const addWorkoutSelectElement = ref(null);
    const addWorkoutSelect = ref(null);
    const categories = ref([]);
    const fetchWorkouts = async () => {
      try {
        const response = await fetch("/app/workouts", {
          headers: new Headers({ "Content-Type": "application/json" }),
        });
        const { data } = await response.json();
        if (response.status != 200) return;

        return data?.workouts ?? [];
      } catch (error) {
        console.error(error);
        return [];
      }
    };
    const initSelect = async () => {
      const workouts = await fetchWorkouts();
      const workoutOptions = workouts.map((w) => ({
        value: w.id,
        label: w.name,
        id: w.id,
        customProperties: w,
      }));
      addWorkoutSelect.value = new Choices(addWorkoutSelectElement.value, {
        allowHTML: false,
      });
      addWorkoutSelect.value.setChoices(workoutOptions);
    };

    const onSubmitNew = async () => {
      try {
        const workouts = addWorkoutSelect.value
          .getValue()
          .map((w) => w.customProperties);
        setValues({ ...values, workouts: workouts });
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
            setErrors({ ...data?.errors });
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
        addWorkoutSelect.value.removeActiveItems();
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
    const deleteCategory = async (id) => {
      const url = `/app/workouts/categories/${id}`;
      const response = await fetch(url, {
        method: "DELETE",
        headers: new Headers({
          "Content-Type": "application/json",
          "X-CSRF-Token": window.csrf,
        }),
      });
      if (response.status === 200) {
        fetchCategories();
        swal.fire(
          "Delete workout category",
          "Category has been deleted.",
          "success"
        );
      }
    };

    const initDelete = async (categoryId) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, delete it",
        title: "Delete workout category",
        text: "Are you sure you want to delete category?",
        confirmButtonColor: "#d9534f",
        cancelButtonText: "I don't want to delete the category",
        icon: "warning",
      });
      if (result.isConfirmed) {
        deleteCategory(categoryId);
      }
    };
    onMounted(() => {
      fetchCategories();
      $("#addCategoryModal").on("hidden.bs.modal", function () {
        resetForm();
        addWorkoutSelect.value.removeActiveItems();
      });
      $("#editCategoryModal").on("hidden.bs.modal", function () {
        resetForm();
      });
      initSelect();
    });
    const name = defineInputBinds("name", { validateOnChange: true });
    // defineInputBinds("workouts", { validateOnChange: true });
    return {
      name,
      errors,
      categories,
      onSubmitNew,
      initEdit,
      onSubmitUpdate,
      initDelete,
      addWorkoutSelectElement,
    };
  },
}).mount("#WorkoutCategoryPage");
