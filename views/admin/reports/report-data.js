import { createApp, onMounted, ref } from "vue";
import ApexChart from "vue3-apexcharts";
createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  components: {
    ApexChart,
  },
  setup() {
    const reportId = window.reportId;
    const breakdownSeries = ref([]);
    const INITIAL_DATA = {
      id: 0,
      clients: 0,
      startDate: "",
      endDate: "",
      coaches: 0,
      members: 0,
      inventoryItems: 0,
      reservations: 0,
      membershipRequests: 0,
      walkIn: null,
      packageRequests: 0,
      earnings: 0,
      earningsBreakdown: {
        walkIn: 0,
        package: 0,
        membership: 0,
      },
      preparedBy: "0",
    };
    const reportData = ref({ ...INITIAL_DATA });
    const fetchReportData = async () => {
      const response = await fetch(`/app/reports/${reportId}`, {
        headers: new Headers({
          "Content-Type": "application/json",
          "Cache-Control": "no-cache",
        }),
      });
      if (response.status === 200) {
        const { data } = await response.json();
        reportData.value = data?.reportData ?? { ...INITIAL_DATA };
        breakdownSeries.value = [
          reportData.value.earningsBreakdown.membership,
          reportData.value.earningsBreakdown.package,
          reportData.value.earningsBreakdown.walkIn,
        ];
      }
    };
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
      },
      labels: ["Membership", "Package", "Walk-In"],
      legend: {
        position: "bottom",
        formatter: (label, data) => {
          return `${label} - ${toMoney(
            breakdownSeries.value[data.seriesIndex]
          )}`;
        },
      },
    };
    const toMoney = (money) => {
      if (!money) return 0;
      return money.toLocaleString(undefined, {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
      });
    };
    onMounted(() => {
      fetchReportData();
    });

    return {
      toMoney,
      reportData,
      pieChartOptions,
      breakdownSeries,
    };
  },
}).mount("#ReportData");
