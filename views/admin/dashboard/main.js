import { createApp, onMounted, ref } from "vue";

createApp({
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
      annualEarnins: 0,
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
        console.log(data);
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
    return {
      dashboardData,
      toMoney,
    };
  },
}).mount("#Dashboard");
