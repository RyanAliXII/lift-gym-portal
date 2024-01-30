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
    const walkInSeries = ref([]);
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
      walkIns: [],
      clientLogs: [],
      packageRequests: 0,
      earnings: 0,
      earningsBreakdown: {
        walkIn: 0,
        package: 0,
        membership: 0,
      },
      preparedBy: "0",
    };

    const barChartOptions = {
      chart: {
        type: "bar",
        toolbar: { show: false },
        animations: {
          enabled: false,
        },
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
        const walkIns =
          reportData.value?.walkIns?.map((walkIn) => ({
            x: formatDate(walkIn.date),
            y: walkIn.total,
          })) ?? [];
        walkInSeries.value = [{ name: "Walk-Ins", data: walkIns }];
      }
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
    const pieChartOptions = {
      chart: {
        type: "pie",
        animations: {
          enabled: false,
        },
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
    const toReadableDatetime = (d) => {
      if (!d) return "";
      const dt = new Date(d);
      try {
        return dt.toLocaleDateString(undefined, {
          month: "long",
          year: "numeric",
          day: "2-digit",
          hour: "2-digit",
          minute: "2-digit",
          hour12: true,
        });
      } catch (error) {
        return "";
      }
    };
    const toReadableDate = (d) => {
      if (!d) return "";
      const dt = new Date(d);
      try {
        return dt.toLocaleDateString(undefined, {
          month: "long",
          year: "numeric",
          day: "2-digit",
        });
      } catch (error) {
        return "";
      }
    };
    onMounted(() => {
      fetchReportData();
    });

    return {
      toMoney,
      reportData,
      pieChartOptions,
      breakdownSeries,
      barChartOptions,
      walkInSeries,
      toReadableDatetime,
      toReadableDate,
      formatDate,
    };
  },
}).mount("#ReportData");
