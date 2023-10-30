package model

import validation "github.com/go-ozzo/ozzo-validation"

type Reservation struct {
	Id int `json:"id" db:"id"`
	DateSlotId int `json:"dateSlotId" db:"date_slot_id"`
	TimeSlotId int `json:"timeSlotId" db:"time_slot_id"`
	ClientId int `json:"clientId" db:"client_id"`
	Model
}

func (m  Reservation) Validate( ) (map[string]string, error) {
	return m.Model.Validate(&m, 
		validation.Field(&m.DateSlotId, validation.Required.Error("Date is required.")),
		validation.Field(&m.TimeSlotId, validation.Required.Error("Time is required.")),
	)
}