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
    let addWorkoutFileUploader = ref(null);
    const { defineInputBinds, errors, values, setErrors, resetForm } = useForm({
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
          headers: new Headers({ "Content-Type": "application/json" }),
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
      addWorkoutFileUploaderGroup.value.appendChild(
        addWorkoutFileUploader.value.element
      );
      fetchWorkouts();
    });
    return {
      name,
      description,
      errors,
      workouts,
      selectedWorkout,
      initView,
      isSubmitting,
      addWorkoutFileUploaderGroup,
      onSubmitNew,
    };
  },
}).mount("#WorkoutPage");
