import { useForm } from "vee-validate";
import { createApp, onMounted, ref } from "vue";
import * as FilePond from "filepond";
import FilePondPluginFileValidateType from "filepond-plugin-file-validate-type";
import swal from "sweetalert2";
FilePond.registerPlugin(FilePondPluginFileValidateType);
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    let addWorkoutFileUploaderGroup = ref(null);
    let editWorkoutFileUploaderGroup = ref(null);
    let addWorkoutFileUploader = ref(null);
    let editWorkoutFileUploader = ref(null);
    const {
      defineInputBinds,
      errors,
      values,
      setErrors,
      resetForm,
      setValues,
    } = useForm({
      initialValues: {
        name: "",
        description: "",
      },
    });
    const isSubmitting = ref(false);
    const workouts = ref([]);
    const selectedWorkout = ref({
      name: "",
      description: "",
      imagePath: "",
      imageSrc: "",
    });
    const fetchWorkouts = async () => {
      try {
        const response = await fetch("/app/workouts", {
          headers: new Headers({
            "Content-Type": "application/json",
            "Cache-Control": "no-cache",
          }),
        });
        const { data } = await response.json();
        if (response.status != 200) return;

        workouts.value = data?.workouts ?? [];
      } catch (error) {
        console.error(error);
      }
    };
    const onSubmitNew = async () => {
      try {
        isSubmitting.value = true;
        if (addWorkoutFileUploader.value.getFiles().length === 0) {
          setErrors({ file: "Animated Image is required." });
          return;
        }
        const fpFile = addWorkoutFileUploader.value.getFile(0);

        const formData = new FormData();
        formData.append("name", values.name);
        formData.append("description", values.description);
        formData.append("file", fpFile.file);

        const response = await fetch("/app/workouts", {
          method: "POST",
          headers: new Headers({ "X-CSRF-Token": window.csrf }),
          body: formData,
        });
        const { data } = await response.json();
        if (response.status >= 400) {
          if (data?.errors) {
            setErrors(data.errors);
          }
          return;
        }
        resetForm();
        addWorkoutFileUploader.value.removeFiles();
        $("#addWorkoutModal").modal("hide");
        fetchWorkouts();
        swal.fire("New Workout", "New workout has been added.", "success");
      } catch (error) {
        console.error(error);
      } finally {
        isSubmitting.value = false;
      }
    };
    const onSubmitUpdate = async () => {
      try {
        isSubmitting.value = true;
        if (editWorkoutFileUploader.value.getFiles().length === 0) {
          setErrors({ file: "Animated Image is required." });
          return;
        }
        const fpFile = editWorkoutFileUploader.value.getFile(0);

        const formData = new FormData();
        formData.append("name", values.name);
        formData.append("description", values.description);
        formData.append("file", fpFile.file);

        const response = await fetch(`/app/workouts/${values.id}`, {
          method: "PUT",
          headers: new Headers({ "X-CSRF-Token": window.csrf }),
          body: formData,
        });
        const { data } = await response.json();
        if (response.status >= 400) {
          if (data?.errors) {
            setErrors(data.errors);
          }
          return;
        }
        resetForm();
        editWorkoutFileUploader.value.removeFiles();
        $("#editWorkoutModal").modal("hide");
        fetchWorkouts();
        swal.fire("Update Workout", "Workout has been updated.", "success");
      } catch (error) {
        console.error(error);
      } finally {
        isSubmitting.value = false;
      }
    };
    const name = defineInputBinds("name", { validateOnChange: true });
    const description = defineInputBinds("description", {
      validateOnChange: true,
    });
    const initView = (workout) => {
      selectedWorkout.value = {
        ...workout,
        imageSrc: `${window.publicURL}/${workout.imagePath}`,
      };
      $("#viewWorkoutModal").modal("show");
    };
    const deleteWorkout = async (id) => {
      const url = `/app/workouts/${id}`;
      const response = await fetch(url, {
        method: "DELETE",
        headers: new Headers({
          "Content-Type": "application/json",
          "X-CSRF-Token": window.csrf,
        }),
      });
      if (response.status === 200) {
        fetchWorkouts();
        swal.fire("Delete workout", "Workout has been deleted.", "success");
      }
    };

    const initDelete = async (workoutId) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, delete it",
        title: "Delete workout",
        text: "Are you sure you want to delete workout?",
        confirmButtonColor: "#d9534f",
        cancelButtonText: "I don't want to delete the workout",
        icon: "warning",
      });
      if (result.isConfirmed) {
        deleteWorkout(workoutId);
      }
    };
    onMounted(() => {
      addWorkoutFileUploader.value = FilePond.create({
        multiple: false,
        acceptedFileTypes: [
          "image/png",
          "image/jpeg",
          "image/jpg",
          "image/gif",
        ],
      });
      editWorkoutFileUploader.value = FilePond.create({
        multiple: false,
        acceptedFileTypes: [
          "image/png",
          "image/jpeg",
          "image/jpg",
          "image/gif",
        ],
      });

      addWorkoutFileUploaderGroup.value.appendChild(
        addWorkoutFileUploader.value.element
      );
      editWorkoutFileUploaderGroup.value.appendChild(
        editWorkoutFileUploader.value.element
      );
      $("#addWorkoutModal").on("hidden.bs.modal", function () {
        addWorkoutFileUploader.value.removeFiles();
        resetForm();
      });
      $("#editWorkoutModal").on("hidden.bs.modal", function () {
        editWorkoutFileUploader.value.removeFiles();
        resetForm();
      });
      fetchWorkouts();
    });

    const initEdit = (workout) => {
      setValues(workout);
      addWorkoutFileUploader.value.removeFiles();
      editWorkoutFileUploader.value.removeFiles();
      const imageSrc = `${window.publicURL}/${workout.imagePath}`;
      editWorkoutFileUploader.value.addFile(imageSrc);
      $("#editWorkoutModal").modal("show");
    };
    return {
      name,
      description,
      errors,
      workouts,
      selectedWorkout,
      initView,
      isSubmitting,
      onSubmitUpdate,
      initEdit,
      initDelete,
      addWorkoutFileUploaderGroup,
      editWorkoutFileUploaderGroup,
      onSubmitNew,
    };
  },
}).mount("#WorkoutPage");
