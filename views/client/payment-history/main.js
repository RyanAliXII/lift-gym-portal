import { createApp, onMounted, ref } from "vue";
import DataTable from "datatables.net-vue3";
import DataTableCore from "datatables.net";
import { formatInTimeZone, zonedTimeToUtc } from "date-fns-tz";
DataTable.use(DataTableCore);
import "datatables.net-dt/css/jquery.dataTables.min.css";
import { format } from "date-fns";
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
        title: "Created At",
        data: "createdAt",
        render: (value) => {
          return formatDate(value);
        },
      },
      {
        title: "Description",
        data: "description",
        render: (value) => {
          return `<span class='font-weight-bold'>${value}</span>`;
        },
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
        const response = await fetch("/clients/payment-history", {
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
    onMounted(() => {
      fetchData();
    });

    return {
      payments,
      columns,
      dtConfig,
    };
  },
}).mount("#PaymentHistory");
