import { useForm } from "vee-validate";
import { createApp, onMounted, ref } from "vue";
import * as FilePond from "filepond";
import FilePondPluginFileValidateType from "filepond-plugin-file-validate-type";
FilePond.registerPlugin(FilePondPluginFileValidateType);
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    let addWorkoutFileUploaderGroup = ref(null);
    let addWorkoutFileUploader = ref(null);
    const { defineInputBinds, errors } = useForm({
      initialValues: {
        name: "",
        description: "",
      },
    });
    const onSubmitNew = async () => {};
    const name = defineInputBinds("name", { validateOnChange: true });
    const description = defineInputBinds("description", {
      validateOnChange: true,
    });
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
    });
    return {
      name,
      description,
      errors,
      addWorkoutFileUploaderGroup,
      onSubmitNew,
    };
  },
}).mount("#WorkoutPage");
