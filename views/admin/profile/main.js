import { createApp, onMounted, ref } from "vue";
import * as FilePond from "filepond";
import FilePondPluginFileValidateType from "filepond-plugin-file-validate-type";
import swal from "sweetalert2";
import { tr } from "date-fns/locale";
FilePond.registerPlugin(FilePondPluginFileValidateType);

createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const avatarUpload = ref(null);
    const errors = ref({});
    const onSubmit = async (event) => {
      errors.value = {};
      const form = new FormData(event.target);
      try {
        const response = await fetch("/app/profile/password", {
          method: "PATCH",
          body: form,
          headers: new Headers({
            "X-CSRF-Token": window.csrf,
          }),
        });
        const { data } = await response.json();
        if (response.status >= 400) {
          if (data?.errors) {
            errors.value = data?.errors;
          }
          return;
        }
        swal.fire(
          "Change Password",
          "Your password has been changed.",
          "success"
        );
        event.target.reset();
      } catch (err) {
        console.error(err);
      }
    };
    onMounted(() => {
      const fp = FilePond.create(avatarUpload.value, {
        labelIdle:
          "Drag and drop your avatar here, or click this to update avatar.",
        maxFiles: 1,
        allowMultiple: false,
        allowFileTypeValidation: true,
        acceptedFileTypes: ["image/png", "image/jpeg", "image/jpg"],
        server: {
          url: "/app/profile/avatar",
          headers: {
            "X-CSRF-Token": window.csrf,
          },
        },
      });
    });
    return {
      onSubmit,
      errors,
      avatarUpload,
    };
  },
}).mount("#ProfilePage");
