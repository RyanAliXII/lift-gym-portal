export const fetchMembers = async () => {
  try {
    const response = await fetch("/app/members", {
      headers: new Headers({ "content-type": "application/json" }),
    });
    const { data } = await response.json();
    return data?.members ?? [];
  } catch (error) {
    console.error(error);
    return;
  }
};
export const fetchClients = async () => {
  try {
    const response = await fetch("/app/clients?type=unsubscribed", {
      headers: new Headers({
        "content-type": "application/json",
        "Cache-Control": "no-cache",
      }),
    });
    const { data } = await response.json();
    return data?.clients;
  } catch (error) {
    console.error(error);
    return [];
  }
};
export const fetchMembershipPlans = async () => {
  try {
    const response = await fetch("/app/memberships", {
      headers: new Headers({
        "content-type": "application/json",
        "Cache-Control": "no-cache",
      }),
    });
    const { data } = await response.json();
    return data?.membershipPlans;
  } catch (error) {
    console.error(error);
    return [];
  }
};
export const subscribe = async (form = {}, onSuccess = () => {}) => {
  try {
    const response = await fetch("/app/members", {
      method: "POST",
      body: JSON.stringify(form),
      headers: new Headers({
        "content-type": "application/json",
        "X-CSRF-Token": window.csrf,
      }),
    });
    if (response.status === 200) {
      onSuccess();
    }
  } catch (error) {
    console.error(error);
  } finally {
    $("#subscribeClientModal").modal("hide");
  }
};

export const cancelSubscription = async (id = 0, onSuccess = () => {}) => {
  try {
    const response = await fetch(`/app/subscriptions/${id}`, {
      method: "DELETE",
      headers: new Headers({
        "content-type": "application/json",
        "X-CSRF-Token": window.csrf,
      }),
    });
    if (response.status === 200) {
      onSuccess();
    }
  } catch (error) {
    console.error(error);
  }
};
