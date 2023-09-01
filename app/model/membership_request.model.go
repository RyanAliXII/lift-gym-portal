package model

import validation "github.com/go-ozzo/ozzo-validation"

type MembershipRequest struct {
	Id               int    `json:"id"`
	ClientId         int    `json:"clientId"`
	MembershipPlanId int    `json:"membershipPlanId"`
	StatusId         int    `json:"statusId"`
	Status           string `json:"status"`
	Client			 ClientJSON	 `json:"client"`
	MembershipPlan 	 MembershipPlanJSON `json:"membershipPlan"`
	Model
}

func (r MembershipRequest) Validate() (error, map[string]string) {
	return r.ValidationRules(&r, 
		validation.Field(&r.ClientId,validation.Required, validation.Min(1)),
		validation.Field(&r.MembershipPlanId,validation.Required, validation.Min(1)),
		validation.Field(&r.StatusId,validation.Required, validation.Min(1)),
	)
}