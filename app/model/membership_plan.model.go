package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type MembershipPlan struct {
	Id          int     `json:"id" db:"id"`
	Price       float64 `json:"price" db:"price"`
	Description string  `json:"description" db:"description"`
	Months      int     `json:"months" db:"months"`
	Model
}

func (m MembershipPlan) Validate() (error, map[string]string){
	return m.ValidationRules(&m, 
		 validation.Field(&m.Price,
		 validation.Required, validation.Min(float64(1))), 
		 validation.Field(&m.Description, validation.Required, validation.Length(1, 300)),
		 validation.Field(&m.Months, validation.Required, validation.Min(1)))
}