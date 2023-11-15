import { createApp, onMounted, ref } from "vue";
import swal from "sweetalert2";
import DataTable from "datatables.net-vue3";
import DataTablesCore from "datatables.net";
import "datatables.net-dt/css/jquery.dataTables.min.css";
DataTable.use(DataTablesCore);
createApp({
  components: {
    DataTable,
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const clients = ref([]);
    const table = ref(null);
    let dt;
    const tableConfig = {
      lengthMenu: [15],
      lengthChange: false,
      dom: "lrtip",
    };
    const columns = [
      {
        data: "publicId",
        title: "ID",
        render: (value) => {
          return `<span class='font-weight-bold'>
            ${value}
          </span>`;
        },
      },
      {
        data: "givenName",
        title: "Given name",
      },
      {
        data: "middleName",
        title: "Middle name",
      },
      {
        data: "surname",
        title: "Surname",
      },
      {
        data: null,
        render: (value, event, row) => {
          let buttons = ``;

          if (window.hasEditPermission) {
            buttons =
              buttons +
              `   <a
          href="/app/clients/${row.id}"
          class="btn btn-outline-primary"
          >
          <i class="fas fa-edit"></i>
         </a>`;
          }
          if (window.hasDeletePermission) {
            buttons =
              buttons +
              `  <button
          class="btn btn-outline-danger delete-client"
          data-toggle="tooltip"
          title="Delete client"
          data-id = ${row.id}
       >
        <i class="fas fa-trash"></i>
       </button>`;
          }

          if (window.hasEditPermission && !row.isVerified) {
            buttons =
              buttons +
              `<button  class="btn btn-outline-success verify-client"    
               data-toggle="tooltip"
              title="Verify account"
              data-id = ${row.id}
              >
              <i class="fa fa-check" aria-hidden="true"></i></button>`;
          }
          if (window.hasEditPermission && row.isVerified) {
            buttons =
              buttons +
              `<button  class="btn btn-outline-warning remove-verify"    
               data-toggle="tooltip"
              title="Remove Verification"
              data-id = ${row.id}
              >
              <i class="fa fa-times" aria-hidden="true"></i></button>`;
          }

          return `<div class='d-flex' style='gap: 5px;'>
            ${buttons}  
          </div>`;
        },
      },
    ];

    const fetchClients = async () => {
      const response = await fetch("/app/clients", {
        headers: new Headers({
          "Content-Type": "application/json",
          "Cache-Control": "no-cache",
        }),
      });
      const { data } = await response.json();
      clients.value = data?.clients ?? [];
    };

    const handleSearch = (event) => {
      const query = event.target.value;
      dt.table().search(query).draw();
    };
    onMounted(() => {
      fetchClients();
      dt = table.value.dt;
      $(dt.table().body()).on(
        "click",
        "button.delete-client",
        async (event) => {
          const id = event.currentTarget.getAttribute("data-id");
          initDelete(id);
        }
      );

      $(dt.table().body()).on(
        "click",
        "button.verify-client",
        async (event) => {
          const id = event.currentTarget.getAttribute("data-id");
          initVerification(id);
        }
      );
      $(dt.table().body()).on(
        "click",
        "button.remove-verify",
        async (event) => {
          const id = event.currentTarget.getAttribute("data-id");
          initUnverification(id);
        }
      );
    });
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
          fetchClients();
        }
      } catch (error) {
        console.error(error);
      }
    };

    const verifyClient = async (id) => {
      try {
        const response = await fetch(`/app/clients/${id}/verification`, {
          method: "PATCH",
          headers: new Headers({ "X-CSRF-Token": window.csrf }),
        });
        if (response.status === 200) {
          await swal.fire(
            "Client Verification",
            "Client has been verified.",
            "success"
          );
          fetchClients();
        }
      } catch (error) {
        console.error(error);
      }
    };
    const unverifyClient = async (id) => {
      try {
        const response = await fetch(`/app/clients/${id}/unverification`, {
          method: "PATCH",
          headers: new Headers({ "X-CSRF-Token": window.csrf }),
        });
        if (response.status === 200) {
          await swal.fire(
            "Client Verification",
            "Client verification has been removed.",
            "success"
          );
          fetchClients();
        }
      } catch (error) {
        console.error(error);
      }
    };
    const initDelete = async (id) => {
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
    const initVerification = async (id) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, verify it.",
        title: "Verify Client",
        text: "Are you sure you want to verify client?",
        confirmButtonColor: "#295ad6",
        cancelButtonText: "I don't want to verify client.",
        icon: "question",
      });
      if (!result.isConfirmed) return;
      verifyClient(id);
    };
    const initUnverification = async (id) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, remove it",
        title: "Remove verification",
        text: "Are you sure you want to remove client verification?",
        confirmButtonColor: "#d9534f",
        cancelButtonText: "I don't want to remove it",
        icon: "warning",
      });
      if (!result.isConfirmed) return;
      unverifyClient(id);
    };
    return {
      table,
      tableConfig,
      columns,
      clients,
      initDelete,
      handleSearch,
    };
  },
}).mount("#ClientPage");
