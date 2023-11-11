import { computed, createApp, onMounted, ref, watch } from "vue";
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
    const INITIAL_DASHBOARD_DATA = {
      clients: 0,
      members: 0,
      annualEarnings: 0,
      monthlyEarnings: 0,
      weeklyEarnings: 0,
      annualEarningsBreakdown: INITAL_BREAKDOWN_DATA,
      monthlyEarningsBreakdown: INITAL_BREAKDOWN_DATA,
      weeklyEarningsBreakdown: INITAL_BREAKDOWN_DATA,
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

    const annualBreakdownSeries = computed(() => {
      return [
        dashboardData.value.annualEarningsBreakdown.membership,
        dashboardData.value.annualEarningsBreakdown.package,
        dashboardData.value.annualEarningsBreakdown.walkIn,
      ];
    });
    const monthlyBreakdownSeries = computed(() => {
      return [
        dashboardData.value.monthlyEarningsBreakdown.membership,
        dashboardData.value.monthlyEarningsBreakdown.package,
        dashboardData.value.monthlyEarningsBreakdown.walkIn,
      ];
    });
    const weeklyBreakdownSeries = computed(() => {
      return [
        dashboardData.value.weeklyEarningsBreakdown.membership,
        dashboardData.value.weeklyEarningsBreakdown.package,
        dashboardData.value.weeklyEarningsBreakdown.walkIn,
      ];
    });
    const pieChartOptions = {
      chart: {
        type: "pie",
      },
      tooltip: {
        y: {
          formatter: function (value) {
            return `${toMoney(value)}`;
          },
        },
        legend: {
          position: "bottom",
        },
      },
      labels: ["Membership", "Package", "Walk In"],
    };

    return {
      dashboardData,
      toMoney,
      annualBreakdownSeries,
      monthlyBreakdownSeries,
      weeklyBreakdownSeries,
      pieChartOptions,
    };
  },
}).mount("#Dashboard");
