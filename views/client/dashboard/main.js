const { createApp, onMounted, ref } = require("vue");
import ApexChart from "vue3-apexcharts";
import DataTable from "datatables.net-vue3";
import DataTablesCore from "datatables.net";
import "datatables.net-dt/css/jquery.dataTables.min.css";
DataTable.use(DataTablesCore);
createApp({
  components: {
    ApexChart,
    DataTable,
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup: () => {
    const INITIAL_DATA = {
      reservations: 0,
      packages: 0,
      coachAppointments: 0,
      penalty: 0,
      annualExpendituresBreakdown: {
        walkIn: 0,
        package: 0,
        membership: 0,
      },
      monthlyExpendituresBreakdown: {
        walkIn: 0,
        package: 0,
        membership: 0,
      },
      weeklyExpendituresBreakdown: {
        walkIn: 0,
        package: 0,
        membership: 0,
      },
      walkIns: [],
    };

    const table = ref(null);
    const columns = [
      {
        title: "Created At",
        data: "createdAt",
        render: (value) => {
          return `<span>${formatDate(value)}</span>`;
        },
      },
      {
        title: "Member(Yes/No)",
        data: "isMember",
        render: (value) => {
          return value ? "Yes" : "No";
        },
      },

      {
        title: "Amount Paid",
        data: "amountPaid",
        render: (value) => toMoney(value),
      },
    ];
    const tableConfig = {
      lengthMenu: [5],
      lengthChange: false,
      dom: "lrtip",
      order: [],
    };
    const formatDate = (date) => {
      return new Date(date).toLocaleDateString("en-US", {
        month: "long",
        day: "2-digit",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
      });
    };
    const expendituresNavs = ref(null);

    const expendituresOverviewSeries = ref([]);
    const expendituresOverviewActiveTab = ref("Weekly");

    const switchExpendituresOverviewTab = (tab) => {
      expendituresOverviewActiveTab.value = tab;
      const navTabId = `#earnings${tab}`;
      const navTab = expendituresNavs.value.querySelector(navTabId);
      const activeNav = expendituresNavs.value.querySelector(".btn-primary");
      activeNav.classList.remove("btn-primary");
      activeNav.classList.add("btn-outline-secondary");

      navTab.classList.add("btn-primary");
      navTab.classList.remove("btn-outline-secondary");
      if (tab === "Monthly") {
        expendituresOverviewSeries.value = [
          dashboardData.value.monthlyExpendituresBreakdown.membership,
          dashboardData.value.monthlyExpendituresBreakdown.package,
          dashboardData.value.monthlyExpendituresBreakdown.walkIn,
        ];
        return;
      }
      if (tab === "Annual") {
        expendituresOverviewSeries.value = [
          dashboardData.value.annualExpendituresBreakdown.membership,
          dashboardData.value.annualExpendituresBreakdown.package,
          dashboardData.value.annualExpendituresBreakdown.walkIn,
        ];
        return;
      }

      expendituresOverviewSeries.value = [
        dashboardData.value.weeklyExpendituresBreakdown.membership,
        dashboardData.value.weeklyExpendituresBreakdown.package,
        dashboardData.value.weeklyExpendituresBreakdown.walkIn,
      ];
    };
    const dashboardData = ref({ ...INITIAL_DATA });
    const fetchDashboardData = async () => {
      const response = await fetch("/clients/dashboard", {
        headers: new Headers({
          "Content-Type": "application/json",
          "Cache-Control": "no-cache",
        }),
      });
      if (response.status === 200) {
        const { data } = await response.json();
        dashboardData.value = data?.dashboardData ?? { ...INITIAL_DATA };

        expendituresOverviewSeries.value = [
          dashboardData.value.weeklyExpendituresBreakdown.membership,
          dashboardData.value.weeklyExpendituresBreakdown.package,
          dashboardData.value.weeklyExpendituresBreakdown.walkIn,
        ];
      }
    };
    const toMoney = (money) => {
      if (!money) return 0;
      return money.toLocaleString(undefined, {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
      });
    };

    const pieChartOptions = {
      chart: {
        type: "pie",
      },
      legend: {
        position: "bottom",
      },
      responsive: [
        {
          breakpoint: 450,
          options: {
            width: 600,
          },
        },
      ],
      tooltip: {
        y: {
          formatter: function (value) {
            return `${toMoney(value)}`;
          },
        },
      },
      labels: ["Membership", "Package", "Walk In"],
    };

    onMounted(() => {
      fetchDashboardData();
    });

    return {
      dashboardData,
      expendituresNavs,
      pieChartOptions,
      switchExpendituresOverviewTab,
      expendituresOverviewActiveTab,
      expendituresOverviewSeries,
      table,
      columns,
      toMoney,
      tableConfig,
    };
  },
}).mount("#Dashboard");
