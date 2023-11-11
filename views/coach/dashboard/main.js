import { createApp, onMounted, ref } from "vue";
import ApexChart from "vue3-apexcharts";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  components: {
    ApexChart,
  },
  setup: () => {
    const INITIAL_DATA = {
      appointments: 0,
      clients: 0,
      earnings: 0,
      weeklyCoachClients: [],
      monthlyCoachClients: [],
    };
    const coachClientNavs = ref(null);

    const dashboardData = ref({ ...INITIAL_DATA });
    const coachedClientSeries = ref([]);
    const toMoney = (money) => {
      if (!money) return 0;
      return money.toLocaleString(undefined, {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
      });
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
    const coachClientTab = ref("Weekly");
    const switchCoachedClientTab = (tab) => {
      coachClientTab.value = tab;
      const navTabId = `#coachedClient${tab}`;
      const navTab = coachClientNavs.value.querySelector(navTabId);
      const activeNav = coachClientNavs.value.querySelector(".btn-primary");

      activeNav.classList.remove("btn-primary");
      activeNav.classList.add("btn-outline-secondary");

      navTab.classList.add("btn-primary");
      navTab.classList.remove("btn-outline-secondary");

      if (tab === "Monthly") {
        const walkIns = dashboardData.value.monthlyCoachClients.map(
          (walkIn) => ({
            x: formatDate(walkIn.date),
            y: walkIn.total,
          })
        );
        coachedClientSeries.value = [{ name: "Walk-Ins", data: walkIns }];
        return;
      }

      const seriesData = dashboardData.value?.weeklyCoachClients.map(
        (coached) => ({
          x: formatDate(coached.date),
          y: coached.total,
        })
      );

      coachedClientSeries.value = [
        {
          name: "Trained Clients",
          data: seriesData,
        },
      ];
    };
    const fetchCoachesDashboardData = async () => {
      const response = await fetch("/coaches/dashboard", {
        headers: new Headers({
          "Content-Type": "application/json",
          "Cache-Control": "no-cache",
        }),
      });
      if (response.status === 200) {
        const { data } = await response.json();
        dashboardData.value = data?.dashboardData ?? { ...INITIAL_DATA };
        const seriesData = dashboardData.value?.weeklyCoachClients.map(
          (coached) => ({
            x: formatDate(coached.date),
            y: coached.total,
          })
        );

        coachedClientSeries.value = [
          {
            name: "Trained Clients",
            data: seriesData,
          },
        ];
      }
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

    onMounted(() => {
      fetchCoachesDashboardData();
    });
    return {
      dashboardData,
      toMoney,
      barChartOptions,
      coachedClientSeries,
      coachClientNavs,
      switchCoachedClientTab,
      coachClientTab,
    };
  },
}).mount("#Dashboard");
