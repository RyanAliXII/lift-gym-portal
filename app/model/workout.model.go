package model

import validation "github.com/go-ozzo/ozzo-validation"

type WorkoutCateory struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Model
}

func (m WorkoutCateory) Validate() (error, map[string]string) {
	return m.Model.ValidationRules(&m, validation.Field(&m.Name, validation.Required.Error("Name is required."), validation.Length(1, 100).Error("Name must be atleast 1 to 100 characters long.")))
}