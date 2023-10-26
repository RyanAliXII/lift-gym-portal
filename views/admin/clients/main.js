import { createApp } from "vue";
import swal from "sweetalert2";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const initDelete = async (event) => {
      const id = event.target.getAttribute("client-id");
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, delete it",
        title: "Delete Client",
        text: "Are you sure you want to delete client?",
        confirmButtonColor: "#d9534f",
        cancelButtonText: "I don't want to delete this client",
        icon: "warning",
      });
      if (result.isConfirmed) {
      }
    };
    return {
      initDelete,
    };
  },
}).mount("#ClientPage");
