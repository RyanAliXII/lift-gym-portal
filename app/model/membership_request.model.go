package model

import validation "github.com/go-ozzo/ozzo-validation"

type MembershipRequest struct {
	Id               int    `json:"id" db:"id"`
	ClientId         int    `json:"clientId" db:"client_id"`
	MembershipPlanId int    `json:"membershipPlanId" db:"membership_plan_id"`
	StatusId         int    `json:"statusId" db:"status_id"`
	Status           string `json:"status" db:"status"`
	Client			 ClientJSON	 `json:"client" db:"client"`
	Remarks          string `json:"remarks" db:"remarks"`
	MembershipPlan 	 MembershipPlanJSON `json:"membershipPlan" db:"membership_plan"`
	MembershipSnapshot MembershipPlanJSON `json:"membershipSnapshot" db:"membership_plan_snapshot"`
	CreatedAt 		 string `json:"createdAt" db:"created_at"`
	Model
}

func (r MembershipRequest) Validate() (error, map[string]string) {
	return r.ValidationRules(&r, 
		validation.Field(&r.ClientId,validation.Required, validation.Min(1)),
		validation.Field(&r.MembershipPlanId,validation.Required, validation.Min(1)),
		validation.Field(&r.StatusId,validation.Required, validation.Min(1)),
	)
}