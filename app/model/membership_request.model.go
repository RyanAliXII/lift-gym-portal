package model

type MembershipRequest struct {
	Id               int    `json:"id"`
	ClientId         int    `json:"clientId"`
	MembershipPlanId int    `json:"membershipPlanId"`
	StatusId         int    `json:"statusId"`
	Status           string `json:"status"`
}
