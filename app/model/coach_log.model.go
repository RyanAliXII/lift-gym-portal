package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type CoachLog struct {
	Id       int  `json:"id" db:"id"`
	CoachId int  `json:"coachId" db:"coach_id"`
	Coach CoachJSON`json:"coach" db:"coach"`
	CreatedAt string `json:"createdAt,omitempty" db:"created_at"`
	IsLoggedOut bool `json:"isLoggedOut" db:"is_logged_out"`
	LoggedOutAt string `json:"loggedOutAt" db:"logged_out_at"`
	Model
}

func (m CoachLog) Validate() (error, map[string]string) {
	return m.Model.ValidationRules(&m, 
		validation.Field(&m.CoachId, validation.Required.Error("Coach is required"), validation.Min(1).Error("Client is required.")))
}
