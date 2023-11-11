import { createApp, ref, onMounted } from "vue";
import swal from "sweetalert2";
import DataTable from "datatables.net-vue3";
import DataTableCore from "datatables.net";
import "datatables.net-dt/css/jquery.dataTables.min.css";

DataTable.use(DataTableCore);
const fetchMembershipRequests = async () => {
  try {
    const response = await fetch("/app/membership-requests", {
      method: "GET",
      headers: new Headers({
        "Content-Type": "application/json",
        "Cache-Control": "no-cache",
      }),
    });
    const { data } = await response.json();
    return data.membershipRequests ?? [];
  } catch (error) {
    console.error(error);
    return [];
  }
};

const updateStatus = async (
  id,
  status = 0,
  onSuccess = () => {},
  data = new FormData()
) => {
  try {
    const response = await fetch(
      `/app/membership-requests/${id}/status?statusId=${status}`,
      {
        method: "PATCH",
        body: data,
        headers: new Headers({
          "X-CSRF-Token": window.csrf,
        }),
      }
    );
    if (response.status === 200) {
      onSuccess();
    }
  } catch (error) {
    console.error(error);
  }
};
createApp({
  components: {
    DataTable,
  },
  setup() {
    const membershipRequests = ref([]);
    const Status = {
      Pending: 1,
      Approved: 2,
      Received: 3,
      Cancelled: 4,
    };

    const table = ref(null);
    let dt;
    const tableConfig = {
      lengthMenu: [20],
      lengthChange: false,
      dom: "lrtip",
    };
    const columns = [
      {
        title: "Created At",
        data: "createdAt",
        render: (value) => formatDate(value),
      },
      {
        title: "Client ID",
        data: "client.publicId",
      },
      {
        title: "Client",
        data: null,
        render: (value, event, row) =>
          `${row.client.givenName} ${row.client.surname}`,
      },
      {
        title: "Membership Plan",
        data: "membershipSnapshot.description",
      },

      {
        title: "Status",
        data: "statusId",
        render: (statusId, event, row) => {
          let text = ``;
          if (statusId === Status.Pending) {
            text = `
              <div class='text-warning'>
                <div class='font-weight-bold'>${row.status}</div>
                <small>Waiting for approval.</small>
              </div>   
             `;
          }
          if (statusId === Status.Approved) {
            text = `
               <div class='text-primary'>
                <div class='font-weight-bold'>${row.status}</div>
                <small>Membership has been approved. Waiting for payment.</small>
               </div>
            `;
          }
          if (statusId === Status.Received) {
            text = `
            <div class='text-success'>
              <div class='font-weight-bold'>${row.status}</div>
              <small >Membership has been received by client.</small>
            </div>
            `;
          }
          if (statusId === Status.Cancelled) {
            text = `
                <div>${row.status}</div>
              `;
          }
          return text;
        },
      },
      {
        title: "Remarks",
        data: "remarks",
        render: (remarks) => {
          if (remarks.length === 0) return "No Remarks.";
          return remarks;
        },
      },
      {
        title: "",
        data: "statusId",
        render: (statusId, event, row) => {
          let buttons = `<div class="d-flex" style="gap: 5px">`;
          if (statusId === Status.Pending) {
            buttons += `<button class="btn btn-outline-primary approve-btn"data-toggle="tooltip" data-placement="top" title="Approve Request" data-id='${row.id}'>
            <i class="fas fa-thumbs-up"></i>
            </button>`;
          }
          if (statusId === Status.Approved) {
            buttons += `<button class="btn btn-outline-primary receive-btn" data-toggle="tooltip" data-placement="top"  title="Mark as Received: Mark this if payment has been made."  data-id='${row.id}'>
            <i class="fas fa-check"></i>
            </button>`;
          }

          if (statusId != Status.Received && statusId != Status.Cancelled) {
            buttons += `<button class="btn btn-outline-danger cancel-btn" data-toggle="tooltip"data-placement="top"title="Cancel Membership Request" data-id='${row.id}'>
              <i class="fas fa-trash"></i>
            </button>`;
          }
          return buttons + `</div>`;
        },
      },
    ];
    const searchRequests = (event) => {
      const query = event.target.value;
      dt.search(query).draw();
    };
    const initApproval = async (id) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, approve it.",
        title: "Approve Membership Request",
        text: "Are you sure you want to approve the membership request?",
        confirmButtonColor: "#295ad6",
        cancelButtonText: "I don't want to approve the request.",
        icon: "question",
      });
      if (result.isConfirmed) {
        updateStatus(id, Status.Approved, async () => {
          swal.fire(
            "Membership Request Approved",
            "Membership request has been approved.",
            "success"
          );
          membershipRequests.value = await fetchMembershipRequests();
        });
      }
    };
    const initCancellation = async (id) => {
      const { value: text, isConfirmed } = await swal.fire({
        input: "textarea",
        title: "Cancel Membership Request",
        inputLabel: "Cancellation Remarks",
        confirmButtonText: "Proceed to cancellation",
        cancelButtonText: "I don't want to cancel the request.",
        confirmButtonColor: "#d9534f",
        inputPlaceholder:
          "Enter cancellation reason eg. duplicate request etc.",
        inputAttributes: {
          "aria-label": "Type your message here",
        },
        showCancelButton: true,
      });
      if (isConfirmed) {
        const formData = new FormData();
        formData.append("remarks", text);
        updateStatus(
          id,
          Status.Cancelled,
          async () => {
            swal.fire(
              "Membership Request Cancellation",
              "Membership request has been cancelled.",
              "success"
            );
            membershipRequests.value = await fetchMembershipRequests();
          },
          formData
        );
      }
    };
    const initMarkAsReceived = async (id) => {
      const result = await swal.fire({
        showCancelButton: true,
        confirmButtonText: "Yes, mark it as received.",
        title: "Recieve Membership Request",
        text: "Are you sure you want to mark request as received?",
        confirmButtonColor: "#295ad6",
        cancelButtonText: "I don't want to mark the request.",
        icon: "question",
      });
      if (result.isConfirmed) {
        updateStatus(id, Status.Received, async () => {
          swal.fire(
            "Membership Request Receiving",
            "Membership has been received by client.",
            "success"
          );
          membershipRequests.value = await fetchMembershipRequests();
        });
      }
    };

    const init = async () => {
      membershipRequests.value = await fetchMembershipRequests();
    };
    const formatDate = (dateStr) => {
      return new Date(dateStr).toLocaleDateString("en-US", {
        month: "long",
        day: "2-digit",
        year: "numeric",
        hour: "2-digit",
        minute: "numeric",
      });
    };
    onMounted(() => {
      init();
      dt = table.value.dt;
      $(dt.table().body()).on("click", "button.approve-btn", async (event) => {
        const id = event.currentTarget.getAttribute("data-id");
        initApproval(id);
      });
      $(dt.table().body()).on("click", "button.receive-btn", async (event) => {
        const id = event.currentTarget.getAttribute("data-id");
        initMarkAsReceived(id);
      });
      $(dt.table().body()).on("click", "button.cancel-btn", async (event) => {
        const id = event.currentTarget.getAttribute("data-id");
        initCancellation(id);
      });
    });
    return {
      membershipRequests,
      Status,
      initApproval,
      initMarkAsReceived,
      initCancellation,
      formatDate,
      searchRequests,
      table,
      tableConfig,
      columns,
    };
  },
  compilerOptions: { delimiters: ["{", "}"] },
}).mount("#MembershipRequest");
