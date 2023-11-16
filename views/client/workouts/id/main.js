import { createApp, onMounted, ref } from "vue";

createApp({
  compilerOptions: {
    delimiters: ["{", "}"],
  },
  setup() {
    const url = new URL(window.location);
    const workouts = ref([]);
    const selectedWorkout = ref({
      id: 0,
      name: "",
      descripion: "",
      imageUrl: "",
    });
    const fetchWorkouts = async () => {
      const response = await fetch(url.pathname, {
        headers: new Headers({
          "Content-Type": "application/json",
          "Cache-Control": "no-cache",
        }),
      });
      if (response.status === 200) {
        const { data } = await response.json();
        workouts.value = data?.workouts ?? [];
      }
    };
    onMounted(() => {
      fetchWorkouts();
    });
    const initView = (workout) => {
      selectedWorkout.value = {
        name: workout.name,
        description: workout.description,
        imageUrl: `${window.publicURL}/${workout.imagePath}`,
      };

      $("#viewWorkoutModal").modal("show");
    };
    return {
      workouts,
      selectedWorkout,
      initView,
    };
  },
}).mount("#Workouts");
