package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Subscribe struct {
	ClientId         int `json:"clientId"`
	MembershipPlanId int `json:"membershipPlanId"`
	MembershipSnapshotId int `json:"membershipSnapshotId"`
	Model
}
func (m *Subscribe) Validate() (error, map[string]string) {
	return m.Model.ValidationRules(m, validation.Field(&m.ClientId, validation.Required, validation.Min(1)), validation.Field(&m.MembershipPlanId, validation.Required, validation.Min(1)))
}

type Member struct {
	Client
	SubscriptionId int `json:"subscriptionId" db:"subscription_id"`
	ValidUntil string `json:"validUntil" db:"valid_until"`	
	SubscriptionStartDate string `json:"subscriptionStartDate" db:"created_at"`
	MembershipPlan MembershipPlanJSON `json:"membershipPlan" db:"membership_plan"`
	MembershipSnapshot MembershipPlanJSON `json:"membershipSnapshot" db:"membership_plan_snaphot"`
}
