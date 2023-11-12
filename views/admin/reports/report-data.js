import { createApp } from "vue";

createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    console.log("test");
    return {};
  },
}).mount("#ReportData");
