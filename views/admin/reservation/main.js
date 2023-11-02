import { createApp, onMounted, ref } from "vue";
import DataTable from "datatables.net-vue3";
import DataTablesCore from "datatables.net";
import "datatables.net-dt/css/jquery.dataTables.min.css";
import swal from "sweetalert2";
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
    const selectedDateSlotId = ref("0");
    const ReservationStatus = {
      Pending: 1,
      Attended: 2,
      NoShow: 3,
      Cancelled: 4,
    };
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
        data: "status",
        title: "Status",
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
          if (row.statusId === ReservationStatus.Attended) {
            return `
            <div class='d-flex' style='gap:5px;'>
              <button class='btn btn-outline-warning no-show-btn' data-toggle="tooltip" title="Mark client as no show" data-id=${row.id}>
              <i class="fa fa-eye-slash" aria-hidden="true"></i>
              </button>
            </div>`;
          }

          if (row.statusId === ReservationStatus.NoShow) {
            return `
            <div class='d-flex' style='gap:5px;'>
             <button class='btn btn-outline-success attended-btn' data-toggle="tooltip" title="Mark client as attended" data-id=${row.id}>
                <i class="fa fa-check" aria-hidden="true"></i
             </button>
            </div>`;
          }
          if (row.statusId === ReservationStatus.Cancelled) {
            return "";
          }
          return `
            <div class='d-flex' style='gap:5px;'>
                 <button class='btn btn-outline-success attended-btn' data-toggle="tooltip" title="Mark client as attended" data-id=${row.id}>
                   <i class="fa fa-check" aria-hidden="true"></i
                 </button>
                 <button class='btn btn-outline-warning  no-show-btn' data-toggle="tooltip" title="Mark client as no show" data-id=${row.id}>
                  <i class="fa fa-eye-slash" aria-hidden="true"></i>
                 </button>
                 <button class='btn btn-outline-danger cancel-btn' data-toggle="tooltip" title="Cancel Reservation" data-id=${row.id}>
                 <i class="fa fa-times" aria-hidden="true"></i>
                </button>
          </div>
        `;
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
        selectedDateSlotId.value = "0";
        fetchReservations();
        return;
      }
      selectedDateSlotId.value = id;
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

    const updateStatus = async (id, statusId = 1) => {
      const response = await fetch(
        `/app/reservations/${id}/status?statusId=${statusId}`,
        {
          method: "PUT",
          headers: new Headers({
            "Content-Type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
        }
      );
      if (response.status === 200) {
        swal.fire(
          "Reservation Status",
          "Reservation status has been updated.",
          "success"
        );

        if (selectedDateSlotId.value == 0) {
          fetchReservations();
          return;
        }
        fetchReservationsDateSlot(selectedDateSlotId.value);
      }
    };
    const search = (event) => {
      const query = event.target.value;
      dt.search(query).draw();
    };
    onMounted(() => {
      dt = table.value.dt;
      fetchReservations();
      $(dt.table().body()).on("click", "button.attended-btn", async (event) => {
        const btn = event.currentTarget;
        const id = btn.getAttribute("data-id");
        const result = await swal.fire({
          showCancelButton: true,
          confirmButtonText: "Yes, mark it",
          confirmButtonColor: "#121717",
          title: "Mark as attended",
          text: "Are you sure you want to mark the reservation as attended?",
          cancelButtonText: "I don't want to mark the reservation",
          icon: "info",
        });
        if (!result.isConfirmed) return;
        updateStatus(id, ReservationStatus.Attended);
      });
      $(dt.table().body()).on("click", "button.no-show-btn", async (event) => {
        const btn = event.currentTarget;
        const id = btn.getAttribute("data-id");
        const result = await swal.fire({
          showCancelButton: true,
          confirmButtonText: "Yes, mark it",
          confirmButtonColor: "#121717",
          title: "Mark as No Show",
          text: "Are you sure you want to mark the reservation as no-show?",
          cancelButtonText: "I don't want to mark the reservation",
          icon: "info",
        });
        if (!result.isConfirmed) return;
        updateStatus(id, ReservationStatus.NoShow);
      });
      $(dt.table().body()).on("click", "button.cancel-btn", async (event) => {
        const btn = event.currentTarget;
        const id = btn.getAttribute("data-id");

        const { value: text, isConfirmed } = await swal.fire({
          input: "textarea",
          inputLabel: "Remarks",
          title: "Cancellation Remarks",
          confirmButtonText: "Submit",
          inputPlaceholder: "Enter the reason for cancellation.",
          inputAttributes: {
            "aria-label": "Enter the reason for cancellation.",
          },
          showCancelButton: true,
          confirmButtonColor: "#d9534f",
        });

        if (!isConfirmed) return;
        updateStatus(id, ReservationStatus.Cancelled);
      });
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
