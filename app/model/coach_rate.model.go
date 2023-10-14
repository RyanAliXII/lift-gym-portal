package model

import validation "github.com/go-ozzo/ozzo-validation"

type CoachRate struct {
	Id          int     `json:"id"`
	Description string  `json:"description" db:"description"`
	Price       float64 `json:"price" db:"price"`
	CoachId     int `json:"coachId" db:"coach_id"`
	Model
}

func (m CoachRate) Validate() (error, map[string]string) {
	return m.ValidationRules(&m, 
		validation.Field(&m.Description, validation.Required.Error("Description is required."), validation.Length(1, 255).Error("Description should be atleast 1 to 255 characters.")),
		validation.Field(&m.Price, validation.Required.Error("Price is required."), validation.Min(float64(1)).Error("Price cannot be less than or equal to 0")),
	)
}
