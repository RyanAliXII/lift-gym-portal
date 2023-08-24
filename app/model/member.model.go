package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Subscribe struct {
	ClientId         int `json:"clientId"`
	MembershipPlanId int `json:"membershipPlanId"`
	Model
}

func (m *Subscribe) Validate() (error, map[string]string) {
	return m.Model.ValidationRules(m, validation.Field(&m.ClientId, validation.Required, validation.Min(1)), validation.Field(&m.MembershipPlanId, validation.Required, validation.Min(1)))
}