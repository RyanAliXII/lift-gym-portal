import { useForm } from "vee-validate";
import { createApp, onMounted, ref } from "vue";
import swal from "sweetalert2";
import { object } from "yup";
import Choices from "choices.js";

const initalErrors = {
  givenName: "",
  middleName: "",
  surname: "",
  email: "",
  roleId: "",
};

const initialForm = {
  id: 0,
  givenName: "",
  middleName: "",
  surname: "",
  email: "",
  gender: "",
  dateOfBirth: "",
  address: "",
  emergencyContact: "",
  mobileNumber: "",

  roleId: 0,
};
createApp({
  setup() {
    const staffs = ref([]);
    const addSelectRoleElement = ref(null);
    const addRoleSelect = ref(null);
    const editSelectRoleElement = ref(null);
    const editRoleSelect = ref(null);
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
    const dateOfBirth = defineInputBinds("dateOfBirth", {
      validateOnChange: true,
    });
    const gender = defineInputBinds("gender", {
      validateOnChange: true,
    });
    const address = defineInputBinds("address", {
      validateOnChange: true,
    });
    const mobileNumber = defineInputBinds("mobileNumber", {
      validateOnChange: true,
    });
    const emergencyContact = defineInputBinds("emergencyContact", {
      validateOnChange: true,
    });
    const fetchStaffs = async () => {
      try {
        const response = await fetch("/app/staffs", {
          headers: new Headers({
            "Content-Type": "application/json",
            "Cache-Control": "no-cache",
          }),
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
        let roleId = addRoleSelect.value.getValue()?.value ?? 0;
        roleId = parseInt(roleId);
        const response = await fetch("/app/staffs", {
          method: "POST",
          body: JSON.stringify({ ...form, roleId }),
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
        let roleId = editRoleSelect.value.getValue()?.value ?? 0;
        roleId = parseInt(roleId);
        const response = await fetch(`/app/staffs/${form.id}`, {
          method: "PUT",
          body: JSON.stringify({ ...form, roleId }),
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
    const initResetPassword = async () => {
      const result = await swal.fire({
        title: "Reset Password",
        text: "Are you sure you want to reset this account password?",
        icon: "warning",
        showCancelButton: true,
        cancelButtonColor: "#3085d6",
        confirmButtonColor: "#d33",
        confirmButtonText: "Yes, reset it",
        cancelButtonText: "I don't want to reset the password",
      });
      if (result.isConfirmed) {
        resetPassword();
      }
    };
    const resetPassword = async () => {
      try {
        const response = await fetch(`/app/staffs/${form.id}/password`, {
          method: "PATCH",
          headers: new Headers({ "X-CSRF-Token": window.csrf }),
        });
        const { data } = await response.json();
        if (response.status === 200) {
          swal.fire(
            "Password Reset",
            `Password reset successful. The password for the staff account is <strong>${data?.password} </strong>`,
            "success"
          );
          $("#editStaffModal").modal("hide");
          fetchStaffs();
        }
      } catch (error) {
        console.error(error);
      }
    };
    const initEdit = (staff) => {
      setValues({ ...staff });
      editRoleSelect.value.setChoiceByValue(staff.roleId.toString());
      $("#editStaffModal").modal("show");
    };
    onMounted(() => {
      addRoleSelect.value = new Choices(addSelectRoleElement.value, {
        allowHTML: true,
      });
      editRoleSelect.value = new Choices(editSelectRoleElement.value, {
        allowHTML: true,
      });
      addRoleSelect.value.removeActiveItems();
      editRoleSelect.value.removeActiveItems();
      fetchStaffs();
      $("#newStaffModal").on("hidden.bs.modal", () => {
        setErrors({ ...initalErrors });
        setValues({ ...initialForm });
        addRoleSelect.value.removeActiveItems();
      });
      $("#editStaffModal").on("hidden.bs.modal", () => {
        setErrors({ ...initalErrors });
        setValues({ ...initialForm });
        editRoleSelect.value.removeActiveItems();
      });
    });
    return {
      givenName,
      middleName,
      surname,
      email,
      mobileNumber,
      address,
      dateOfBirth,
      gender,
      emergencyContact,
      errors,
      onSubmitNewStaff,
      staffs,
      initEdit,
      onSubmitUpdate,
      initResetPassword,
      addSelectRoleElement,
      editRoleSelect,
      editSelectRoleElement,
    };
  },
  compilerOptions: {
    delimiters: ["{", "}"],
  },
}).mount("#StaffPage");
