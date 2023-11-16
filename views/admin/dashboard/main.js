import { createApp, onMounted, ref } from "vue";
import ApexChart from "vue3-apexcharts";
createApp({
  components: {
    ApexChart,
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const INITAL_BREAKDOWN_DATA = {
      package: 0,
      walkIn: 0,
      membership: 0,
    };
    const earningsOverviewNavTabs = ref(null);
    const walkInNavs = ref(null);
    const earningsOverviewSeries = ref([]);
    const earningsOverviewActiveTab = ref("Weekly");
    const walkInSeries = ref([]);
    const walkInOverTab = ref("Weekly");
    const switchEarningsOverviewTab = (tab) => {
      earningsOverviewActiveTab.value = tab;
      const navTabId = `#earnings${tab}`;
      const navTab = earningsOverviewNavTabs.value.querySelector(navTabId);
      const activeNav =
        earningsOverviewNavTabs.value.querySelector(".btn-primary");
      activeNav.classList.remove("btn-primary");
      activeNav.classList.add("btn-outline-secondary");

      navTab.classList.add("btn-primary");
      navTab.classList.remove("btn-outline-secondary");
      if (earningsOverviewActiveTab.value === "Monthly") {
        earningsOverviewSeries.value = [
          dashboardData.value.monthlyEarningsBreakdown.membership,
          dashboardData.value.monthlyEarningsBreakdown.package,
          dashboardData.value.monthlyEarningsBreakdown.walkIn,
        ];
        return;
      }

      if (earningsOverviewActiveTab.value === "Annual") {
        earningsOverviewSeries.value = [
          dashboardData.value.annualEarningsBreakdown.membership,
          dashboardData.value.annualEarningsBreakdown.package,
          dashboardData.value.annualEarningsBreakdown.walkIn,
        ];
        return;
      }
      earningsOverviewSeries.value = [
        dashboardData.value.weeklyEarningsBreakdown.membership,
        dashboardData.value.weeklyEarningsBreakdown.package,
        dashboardData.value.weeklyEarningsBreakdown.walkIn,
      ];
    };
    const switchWalkInTab = (tab) => {
      walkInOverTab.value = tab;
      const navTabId = `#walkIn${tab}`;
      const navTab = walkInNavs.value.querySelector(navTabId);
      const activeNav = walkInNavs.value.querySelector(".btn-primary");
      activeNav.classList.remove("btn-primary");
      activeNav.classList.add("btn-outline-secondary");

      navTab.classList.add("btn-primary");
      navTab.classList.remove("btn-outline-secondary");
      if (tab === "Monthly") {
        const walkIns = dashboardData.value.monthlyWalkIns.map((walkIn) => ({
          x: formatDate(walkIn.date),
          y: walkIn.total,
        }));
        walkInSeries.value = [{ name: "Walk-Ins", data: walkIns }];
        return;
      }
      earningsOverviewSeries.value = [
        dashboardData.value.weeklyEarningsBreakdown.membership,
        dashboardData.value.weeklyEarningsBreakdown.package,
        dashboardData.value.weeklyEarningsBreakdown.walkIn,
      ];
      const walkIns = dashboardData.value.weeklyWalkIns.map((walkIn) => ({
        x: formatDate(walkIn.date),
        y: walkIn.total,
      }));
      walkInSeries.value = [{ name: "Walk-Ins", data: walkIns }];
    };

    const INITIAL_DASHBOARD_DATA = {
      clients: 0,
      members: 0,
      annualEarnings: 0,
      monthlyEarnings: 0,
      weeklyEarnings: 0,
      weeklyWalkIns: [],
      monthlyWalkIns: [],
      annualEarningsBreakdown: INITAL_BREAKDOWN_DATA,
      monthlyEarningsBreakdown: INITAL_BREAKDOWN_DATA,
      weeklyEarningsBreakdown: INITAL_BREAKDOWN_DATA,
    };
    const formatDate = (date) => {
      if (!date) return "No Date";
      if (date.length === 0) return "No Date";
      return new Date(date).toLocaleDateString("en-US", {
        month: "long",
        day: "2-digit",
        year: "numeric",
      });
    };
    const dashboardData = ref({ ...INITIAL_DASHBOARD_DATA });
    const fetchDashboardData = async () => {
      const response = await fetch("/app/dashboard", {
        headers: new Headers({
          "Content-Type": "application/json",
          "Cache-Control": "no-cache",
        }),
      });
      if (response.status === 200) {
        const { data } = await response.json();
        dashboardData.value = data?.dashboardData ?? {
          ...INITIAL_DASHBOARD_DATA,
        };
        earningsOverviewSeries.value = [
          dashboardData.value.weeklyEarningsBreakdown.membership,
          dashboardData.value.weeklyEarningsBreakdown.package,
          dashboardData.value.weeklyEarningsBreakdown.walkIn,
        ];
        const walkIns = dashboardData.value.weeklyWalkIns.map((walkIn) => ({
          x: formatDate(walkIn.date),
          y: walkIn.total,
        }));
        walkInSeries.value = [{ name: "Walk-Ins", data: walkIns }];
      }
    };
    onMounted(() => {
      fetchDashboardData();
    });

    const toMoney = (money) => {
      if (!money) return 0;
      return money.toLocaleString(undefined, {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
      });
    };

    const barChartOptions = {
      chart: {
        type: "bar",
      },
      yaxis: {
        labels: {
          formatter: function (val) {
            return val.toFixed(0);
          },
        },
      },
      plotOptions: {
        bar: {
          distributed: true,
        },
      },
    };

    const pieChartOptions = {
      chart: {
        type: "pie",
      },
      responsive: [
        {
          breakpoint: 450,
          options: {
            width: 600,
          },
        },
      ],
      legend: {
        position: "bottom",
      },
      tooltip: {
        y: {
          formatter: function (value) {
            return `${toMoney(value)}`;
          },
        },
      },
      labels: ["Membership", "Package", "Walk In"],
    };

    return {
      dashboardData,
      toMoney,
      barChartOptions,
      pieChartOptions,
      walkInSeries,
      earningsOverviewActiveTab,
      walkInOverTab,
      switchEarningsOverviewTab,
      switchWalkInTab,
      earningsOverviewSeries,
      earningsOverviewNavTabs,
      walkInNavs,
    };
  },
}).mount("#Dashboard");
