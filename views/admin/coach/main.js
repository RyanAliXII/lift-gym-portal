import swal from "sweetalert2";
import { createApp } from "vue";

createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup: () => {
    const deleteCoach = async (id) => {
      const response = await fetch(`/app/coaches/${id}`, {
        method: "DELETE",
        headers: new Headers({
          "Content-Type": "application/json",
          "X-CSRF-Token": window.csrf,
        }),
      });
      if (response.status === 200) {
        await swal.fire("Delete Coach", "Coach has been deleted.", "success");
        location.reload();
      }
    };
    const initDelete = async (id) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, delete it",
        title: "Delete Coach",
        text: "Are you sure you want to delete coach?",
        confirmButtonColor: "#d9534f",
        cancelButtonText: "I don't want to delete this coach",
        icon: "warning",
      });
      if (result.isConfirmed) {
        deleteCoach(id);
      }
    };
    return {
      initDelete,
    };
  },
}).mount("#CoachPage");
