import { createApp } from "vue";
import swal from "sweetalert2";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const deleteClient = async (id) => {
      try {
        const response = await fetch(`/app/clients/${id}`, {
          method: "DELETE",
          headers: new Headers({ "X-CSRF-Token": window.csrf }),
        });
        if (response.status === 200) {
          await swal.fire(
            "Client Delete",
            "Client has been deleted.",
            "success"
          );
          location.reload();
        }
      } catch (error) {
        console.error(error);
      }
    };
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
        deleteClient(id);
      }
    };
    return {
      initDelete,
    };
  },
}).mount("#ClientPage");