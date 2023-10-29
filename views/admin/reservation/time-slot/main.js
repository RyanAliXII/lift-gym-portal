import { createApp, onMounted, ref } from "vue";

createApp({
  setup() {
    const form = ref({
      startTime: "",
      endTime: "",
      maxCapacity: 1,
    });
    const errors = ref({});

    const handleFormInput = (event) => {
      let value = event.target.value;
      let name = event.target.name;
      if (event.target.type === "number") {
        value = Number(value);
      }
      form.value[name] = value;
      delete errors.value[name];
    };

    const onSubmit = () => {};
    onMounted(() => {});
    return {
      form,
      handleFormInput,
      onSubmit,
    };
  },
}).mount("#TimeSlot");
