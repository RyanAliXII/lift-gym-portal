import { createApp, ref } from "vue";
import DataTable from "datatables.net-vue3";
import DataTableCore from "datatables.net";
DataTable.use(DataTableCore);
import "datatables.net-dt/css/jquery.dataTables.min.css";
createApp({
  components: {
    DataTable,
  },
  setup() {
    const payments = ref([]);
    const dtConfig = {
      lengthMenu: [20],
      dom: "lrtip",
      lengthChange: false,
    };
    const columns = [
      {
        title: "Amount",
        data: "amount",
      },
      {
        title: "Description",
        data: "description",
      },
    ];

    return {
      payments,
      columns,
      dtConfig,
    };
  },
}).mount("#PaymentHistory");
