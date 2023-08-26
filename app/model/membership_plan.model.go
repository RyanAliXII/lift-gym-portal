package model

import (
	"database/sql/driver"
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation"
)

type MembershipPlan struct {
	Id          int     `json:"id" db:"id"`
	Price       float64 `json:"price" db:"price"`
	Description string  `json:"description" db:"description"`
	Months      int     `json:"months" db:"months"`
	Model
}
type MembershipPlanJSON struct {
	MembershipPlan
}
func (instance *MembershipPlanJSON) Scan(value interface{}) error {
	val, valid := value.([]byte)
	if valid {
		unmarshalErr := json.Unmarshal(val, instance)
		if unmarshalErr != nil {
			*instance = MembershipPlanJSON{}
		}
	} else {
		*instance = MembershipPlanJSON{}
	}
	return nil

}
func (copy MembershipPlanJSON) Value(value interface{}) (driver.Value, error) {
	return copy, nil
}
func (m MembershipPlan) Validate() (error, map[string]string){
	return m.ValidationRules(&m, 
		 validation.Field(&m.Price,
		 validation.Required, validation.Min(float64(1))), 
		 validation.Field(&m.Description, validation.Required, validation.Length(1, 300)),
		 validation.Field(&m.Months, validation.Required, validation.Min(1)))
}