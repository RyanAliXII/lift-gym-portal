import { createApp, onMounted, ref } from "vue";
import DataTable from "datatables.net-vue3";
import DataTablesCore from "datatables.net";
import "datatables.net-dt/css/jquery.dataTables.min.css";
import { ta } from "date-fns/locale";
// import "datatables.net-bs4";
// import "datatables.net-responsive-dt";

DataTable.use(DataTablesCore);
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  components: {
    DataTable,
  },
  setup() {
    const table = ref("");
    let dt;
    const options = {
      lengthMenu: [20],
      lengthChange: false,
    };
    const columns = [
      {
        data: "reservationId",
        title: "Reservation Id",
        render: (v) => {
          return `<span class="font-weight-bold">${v}</span>`;
        },
      },
      {
        data: "client",
        title: "Client",
        render: (v) => {
          return `${v.givenName} ${v.surname}`;
        },
      },
      {
        data: "date",
        title: "Date",
        render: (v) => {
          return formatDate(v);
        },
      },
      { data: "time", title: "Time" },

      {
        data: null,
        render: (value, type, row) => {
          return `<button class='btn btn-outline-success attended-btn' data-toggle="tooltip" title="Mark client as attended" row-id=${row.id}>
          <i class="fa fa-check" aria-hidden="true"></i
        </button>`;
        },
      },
    ];
    const reservations = ref([]);
    const fetchReservations = async () => {
      const response = await fetch("/app/reservations", {
        headers: new Headers({ "Content-Type": "application/json" }),
      });
      const { data } = await response.json();
      if (response.status === 200) {
        reservations.value = data?.reservations ?? [];
      }
    };
    const formatDate = (date) => {
      if (!date) return "";
      return new Date(date).toLocaleDateString("en-US", {
        month: "long",
        day: "2-digit",
        year: "numeric",
      });
    };
    const handleDateSelect = (event) => {
      const id = event.target.value;
      if (id == 0) {
        fetchReservations();
        return;
      }
      fetchReservationsDateSlot(id);
    };
    const fetchReservationsDateSlot = async (id) => {
      const response = await fetch(`/app/reservations/date-slots/${id}`, {
        headers: new Headers({ "Content-Type": "application/json" }),
      });
      const { data } = await response.json();
      if (response.status === 200) {
        reservations.value = data?.reservations ?? [];
      }
    };
    const search = (event) => {
      const query = event.target.value;
      dt.search(query).draw();
    };
    onMounted(() => {
      dt = table.value.dt;
      fetchReservations();
      $(dt.table().body()).on("click", "button.attended-btn", (event) => {});
    });
    return {
      reservations,
      formatDate,
      handleDateSelect,
      columns,
      options,
      table,
      search,
    };
  },
}).mount("#ReservationPage");
