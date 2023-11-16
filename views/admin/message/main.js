import { createApp, onMounted, ref } from "vue";
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
    const table = ref();
    const message = ref({
      name: "",
      email: "",
      message: "",
    });
    let dt;
    const messages = ref([]);
    const columns = [
      {
        title: "Email",
        data: "email",
      },
      {
        title: "Name",
        data: "name",
      },
      {
        title: "",
        data: null,
        render: (value, event, row) => {
          return `<button data-id='${row.id}'  class='btn btn-primary view-message'>View</button>`;
        },
      },
    ];
    const search = () => {};
    const tableConfig = {
      lengthMenu: [15],
      lengthChange: false,
      order: [],
      dom: "lrtip",
    };

    const fetchMessages = async () => {
      try {
        const response = await fetch("/app/messages", {
          headers: new Headers({
            "Content-Type": "application/json",
            "Cache-Control": "no-cache",
          }),
        });
        const { data } = await response.json();
        if (response.status === 200) {
          messages.value = data?.messages ?? [];
        }
      } catch (error) {}
    };

    onMounted(() => {
      dt = table.value.dt;
      fetchMessages();
      $(dt.table().body()).on("click", "button.view-message", async (event) => {
        const id = event.currentTarget.getAttribute("data-id");

        const msg = messages.value.find((m) => m.id == id);

        message.value = msg;
        $("#viewMessageModal").modal("show");
      });
    });
    return {
      table,
      columns,
      messages,
      search,
      message,
      tableConfig,
    };
  },
}).mount("#Messages");
