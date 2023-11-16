package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type StaffLog struct {
	Id       int  `json:"id" db:"id"`
	StaffId int  `json:"staffId" db:"staff_id"`
	Staff StaffJSON `json:"staff" db:"staff"`
	CreatedAt string `json:"createdAt,omitempty" db:"created_at"`
	IsLoggedOut bool `json:"isLoggedOut" db:"is_logged_out"`
	LoggedOutAt string `json:"loggedOutAt" db:"logged_out_at"`
	Model
}

func (m StaffLog) Validate() (error, map[string]string) {
	return m.Model.ValidationRules(&m, 
		validation.Field(&m.StaffId, validation.Required.Error("Coach is required"), validation.Min(1).Error("Client is required.")))
}
