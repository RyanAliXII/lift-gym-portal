package model

import validation "github.com/go-ozzo/ozzo-validation"

type Reservation struct {
	Id int `json:"id" db:"id"`
	DateSlotId int `json:"dateSlotId" db:"date_slot_id"`
	TimeSlotId int `json:"timeSlotId" db:"time_slot_id"`
	ClientId int `json:"clientId" db:"client_id"`
	IsCancelled bool `json:"isCancelled" db:"is_cancelled"`
	HasAttended bool `json:"hasAttended" db:"has_attended"`
	StatusId int `json:"statusId" db:"status_id"`
	Status string `json:"status" db:"status"`
	Date string `json:"date" db:"date"`
	Time string `json:"time" db:"time"`
	ReservationId string `json:"reservationId" db:"reservation_id"`
	Client ClientJSON `json:"client,omitempty" db:"client"`
	Model
}

func (m  Reservation) Validate( ) (map[string]string, error) {
	return m.Model.Validate(&m, 
		validation.Field(&m.DateSlotId, validation.Required.Error("Date is required.")),
		validation.Field(&m.TimeSlotId, validation.Required.Error("Time is required.")),
	)
}