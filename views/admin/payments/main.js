import { createApp, onMounted, ref } from "vue";
import DataTable from "datatables.net-vue3";
import DataTableCore from "datatables.net";
DataTable.use(DataTableCore);
import "datatables.net-dt/css/jquery.dataTables.min.css";

createApp({
  components: {
    DataTable,
  },
  setup() {
    const table = ref(null);
    let dt;
    const payments = ref([]);
    const dtConfig = {
      lengthMenu: [20],
      dom: "lrtip",
      lengthChange: false,
      order: [],
    };
    const columns = [
      {
        title: "Created At",
        data: "createdAt",
        render: (value) => {
          return formatDate(value);
        },
      },
      {
        title: "Client ID",
        data: "client",
        render: (value) => {
          return `<span class='font-weight-bold'>${value.publicId}</span>`;
        },
      },
      {
        title: "Client",
        data: "client",
        render: (value) => {
          return `<span class='font-weight-bold'>${value.givenName} ${value.surname}</span>`;
        },
      },
      {
        title: "Description",
        data: "description",
      },
      {
        title: "Amount",
        data: "amount",
        render: (value) => {
          return toMoney(value);
        },
      },
    ];
    const toMoney = (money) => {
      if (!money) return 0;
      return money.toLocaleString(undefined, {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
      });
    };
    const formatDate = (date) => {
      if (!date) return "";
      return new Date(date).toLocaleTimeString(undefined, {
        timeZone: "Asia/Singapore",
        month: "long",
        day: "2-digit",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
      });
    };
    const fetchData = async () => {
      try {
        const response = await fetch("/app/payments", {
          headers: new Headers({
            "Content-Type": "application/json",
            "Cache-Control": "no-cache",
          }),
        });
        if (response.status === 200) {
          const { data } = await response.json();
          payments.value = data?.payments ?? [];
        }
      } catch (err) {
        console.error(err);
      }
    };
    const searchPayments = (event) => {
      const query = event.target.value;
      dt.search(query).draw();
    };
    onMounted(() => {
      dt = table.value.dt;
      fetchData();
    });

    return {
      payments,
      columns,
      dtConfig,
      table,
      searchPayments,
    };
  },
}).mount("#PaymentHistory");
