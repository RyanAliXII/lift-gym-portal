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
    const tableConfig = {
      lengthMenu: [15],
      lengthChange: false,
      dom: "lrtip",
    };
    const columns = [
      {
        data: "publicId",
        title: "Id",
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
    onMounted(() => {
      fetchClients();
      const dt = table.value.dt;
      $(dt.table().body()).on(
        "click",
        "button.delete-client",
        async (event) => {
          const id = event.currentTarget.getAttribute("data-id");
          initDelete(id);
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

    return {
      table,
      tableConfig,
      columns,
      clients,
      initDelete,
    };
  },
}).mount("#ClientPage");
