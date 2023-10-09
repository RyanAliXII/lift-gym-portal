package model

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Equipment struct {
	Name         string  `json:"name" db:"name"`
	ModelOrMake  string  `json:"model" db:"model"`
	Quantity     int     `json:"quantity" db:"quantity"`
	CostPrice    float64 `json:"costPrice" db:"cost_price"`
	DateReceived string  `json:"dateReceived" db:"date_received"`
	Model
}

func (m *Equipment) Validate() {
	m.Model.ValidationRules(&m, 
	validation.Field(&m.Name, validation.Required.Error("Equipment name is required.",), validation.Length(1, 100).Error("Equipment name should be atleast 1 to 100 characters.")),
    validation.Field(&m.ModelOrMake, validation.Required.Error("Model/make is required.",), validation.Length(1, 100).Error("Model/make should be atleast 1 to 100 characters.")),
	validation.Field(&m.Quantity, validation.Required.Error("Quantity is required."), validation.Min(1).Error("Quantity shoud be atleast 1")),
	validation.Field(&m.CostPrice, validation.Required.Error("Cost price is required."), validation.Min(float64(1)).Error("Cost price should be atleast 1")),
	validation.Field(&m.DateReceived, validation.Required.Error("Date received is required."), validation.By(func(value interface{}) error {
		format := "2006-01-02"
		strDate ,_ := value.(string)
		_, err := time.Parse(format, strDate)
		if err != nil {
			return fmt.Errorf("Date of birth is required.")
		}
		return nil
	})))
	
}