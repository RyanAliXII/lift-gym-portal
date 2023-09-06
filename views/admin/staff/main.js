import { useForm } from "vee-validate";
import { createApp, onMounted, ref } from "vue";
import swal from "sweetalert2";
import { object } from "yup";

const initalErrors = {
  givenName: "",
  middleName: "",
  surname: "",
  email: "",
};

const initialForm = {
  id: 0,
  givenName: "",
  middleName: "",
  surname: "",
  email: "",
};
createApp({
  setup() {
    const staffs = ref([]);
    const {
      defineInputBinds,
      values: form,
      errors,
      setValues,
      setErrors,
    } = useForm({
      initialValues: initialForm,
      validationSchema: object(),
    });

    const givenName = defineInputBinds("givenName", {
      validateOnChange: true,
    });
    const middleName = defineInputBinds("middleName", {
      validateOnChange: true,
    });
    const surname = defineInputBinds("surname", {
      validateOnChange: true,
    });
    const email = defineInputBinds("email", {
      validateOnChange: true,
    });
    const fetchStaffs = async () => {
      try {
        const response = await fetch("/app/staffs", {
          headers: new Headers({ "Content-Type": "application/json" }),
        });
        const { data } = await response.json();
        staffs.value = data?.staffs ?? [];
      } catch (error) {
        console.error(error);
      }
    };
    const onSubmitNewStaff = async () => {
      try {
        setErrors({ ...initalErrors });
        const response = await fetch("/app/staffs", {
          method: "POST",
          body: JSON.stringify(form),
          headers: new Headers({
            "content-type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
        });
        const { data } = await response.json();
        if (response.status === 200) {
          swal.fire(
            "New Staff",
            `Staff has been succcessfully added. <br> The password for the staff account is <strong>${data?.password} </strong>`,
            "success"
          );
          $("#newStaffModal").modal("hide");
          fetchStaffs();
        }
        if (response.status === 400 && data?.errors) {
          setErrors(data.errors);
        }
      } catch (error) {
        console.error(error);
      }
    };
    const onSubmitUpdate = async () => {
      try {
        const response = await fetch(`/app/staffs/${form.id}`, {
          method: "PUT",
          body: JSON.stringify(form),
          headers: new Headers({
            "Content-Type": "application/json",
            "X-CSRF-Token": window.csrf,
          }),
        });
        const { data } = await response.json();
        if (response.status === 200) {
          swal.fire("Staff Update", "Staff has been updated.", "success");
          $("#editStaffModal").modal("hide");
          fetchStaffs();
        }
        if (response.status === 400 && data?.errors) {
          setErrors(data.errors);
        }
      } catch (error) {
        console.error(error);
      }
    };
    const initEdit = (staff) => {
      setValues({ ...staff });
      $("#editStaffModal").modal("show");
    };
    onMounted(() => {
      fetchStaffs();
      $("#newStaffModal").on("hidden.bs.modal", () => {
        setErrors({ ...initalErrors });
        setValues({ ...initialForm });
      });
      $("#editStaffModal").on("hidden.bs.modal", () => {
        setErrors({ ...initalErrors });
        setValues({ ...initialForm });
      });
    });
    return {
      givenName,
      middleName,
      surname,
      email,
      errors,
      onSubmitNewStaff,
      staffs,
      initEdit,
      onSubmitUpdate,
    };
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
}).mount("#StaffPage");
