import { createApp, onMounted, ref } from "vue";
import DataTable from "datatables.net-vue3";
import DataTablesCore from "datatables.net";
import "datatables.net-dt/css/jquery.dataTables.min.css";
DataTable.use(DataTablesCore);
createApp({
  components: {
    DataTable,
  },
  setup() {
    let dt;
    const table = ref(null);
    const tableConfig = {
      lengthMenu: [15],
      lengthChange: false,
      dom: "lrtip",
    };

    const columns = [
      {
        title: "Amount",
        data: "amount",
      },
      {
        title: "Client",
        data: "client",
        render: (value) => {
          return `<span class='font-weight-bold'> ${value.givenName} ${value.surname}</span>`;
        },
      },
      {
        title: "Coach",
        data: "coach",
        render: (value) => {
          return `<span class='font-weight-bold'> ${value.givenName} ${value.surname}</span>`;
        },
      },
      {
        title: "Description",
        data: null,
        render: () => {
          return "Did not show up on 3 consecutive coaching session.";
        },
      },
      {
        title: "Status",
        data: "isSettled",
        render: (value) => {
          if (value) return `<span class='text-success'>Settled</span>`;
          return `<span class='text-danger'>Unsettled</span>`;
        },
      },

      {
        title: "",
        data: "isSettled",
        render: (value, event, row) => {
          if (!value)
            return `<button class='btn btn-success settle-btn' data-id=${row.id}>Settle</button>`;

          return `<button class='btn btn-danger unsettle-btn' data-id=${row.id}>Unsettle</button>`;
        },
      },
    ];
    const penalties = ref([]);
    const fetchPenalties = async () => {
      const response = await fetch("/app/coaching-penalty", {
        headers: new Headers({ "Content-Type": "application/json" }),
      });
      if (response.status === 200) {
        const { data } = await response.json();
        penalties.value = data?.penalties ?? [];
      }
    };
    const handleSearch = (event) => {
      const query = event.target.value;
      dt.table().search(query).draw();
    };
    onMounted(() => {
      fetchPenalties();
      dt = table.value.dt;

      $(dt.table().body()).on("click", "button.settle-btn", async (event) => {
        const id = event.currentTarget.getAttribute("data-id");
        console.log(id);
      });

      $(dt.table().body()).on("click", "button.unsettle-btn", async (event) => {
        const id = event.currentTarget.getAttribute("data-id");
        console.log(id);
      });
    });

    return {
      columns,
      penalties,
      tableConfig,
      handleSearch,
      table,
    };
  },
}).mount("#PenaltyPage");
