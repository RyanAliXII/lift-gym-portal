package model

import (
	"fmt"
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/pkg/status"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Reservation struct {
	Id int `json:"id" db:"id"`
	DateSlotId int `json:"dateSlotId" db:"date_slot_id"`
	TimeSlotId int `json:"timeSlotId" db:"time_slot_id"`
	ClientId int `json:"clientId" db:"client_id"`
	StatusId int `json:"statusId" db:"status_id"`
	Status string `json:"status" db:"status"`
	Date string `json:"date" db:"date"`
	Time string `json:"time" db:"time"`
	ReservationId string `json:"reservationId" db:"reservation_id"`
	Client ClientJSON `json:"client,omitempty" db:"client"`
	Remarks string `json:"remarks" db:"remarks"`
	Model
}
func(r  Reservation)TransformDateStrToTime() (time.Time, error) {
	return time.Parse(time.DateOnly, r.Date)
}

func (m  Reservation) Validate( ) (map[string]string, error) {
	return m.Model.Validate(&m, 
		validation.Field(&m.DateSlotId, validation.Required.Error("Date is required.")),
		validation.Field(&m.TimeSlotId, validation.Required.Error("Time is required.",),  validation.By(validateIfTemporarilyBanned(m.ClientId)) ),
	)
}
func validateIfTemporarilyBanned(clientId int) validation.RuleFunc {
	return func( value interface{}) error {
		db := db.GetConnection()
		reservations := make([]Reservation, 0)
		err := db.Select(&reservations,`SELECT date_slot.date, status_id FROM reservation  
			INNER JOIN date_slot on reservation.date_slot_id = date_slot.id
			where client_id = ? and  date_slot.date < curdate()
			ORDER BY date desc LIMIT 3
		`, clientId)
		if err != nil {
			return fmt.Errorf("you are not allowed to create a new reservation")
		}
		
		const MaxAllowedUnattendedReservations = 3
		if len(reservations) < MaxAllowedUnattendedReservations{
			return nil
		}
		unattendedCount := 0;
		for _, reservation := range reservations{
			if(reservation.StatusId == status.ReservationStatusNoShow || reservation.StatusId == status.ReservationStatusPending){
				unattendedCount = unattendedCount + 1
			}
		}
		OneMonth := (time.Hour * 24) * 30
		if(unattendedCount >= MaxAllowedUnattendedReservations){
			latestReservation, err := reservations[0].TransformDateStrToTime()
			if err != nil {
				return fmt.Errorf("you are not allowed to create a new reservation")
			}
			latestReservationOneMonthFromNow := latestReservation.Add(OneMonth)
			today := time.Now().Truncate(24 * time.Hour)
		
			if(today.Before(latestReservationOneMonthFromNow)){
				return fmt.Errorf("you are not allowed to create a new reservation")
			}
		}
		return nil
	}
}

/* 
	this function check if user has 3 consecutive unattended reservation
	if user has 3 consecutive unattended reservation, apply 30 days ban starting from last reservation date.
*/
func IsTemporarilyBannedFromReservation(clientId int) bool {
	db := db.GetConnection()
	reservations := make([]Reservation, 0)
	err := db.Select(&reservations,`SELECT date_slot.date, status_id FROM reservation  
	INNER JOIN date_slot on reservation.date_slot_id = date_slot.id
	where client_id = ? and  date_slot.date < curdate()
	ORDER BY date desc LIMIT 3
	`, clientId)
	if err != nil {
			return true
	}
		
	const MaxAllowedUnattendedReservations = 3
	if len(reservations) < MaxAllowedUnattendedReservations{
			return true
	}
	unattendedCount := 0;
	for _, reservation := range reservations{
		if(reservation.StatusId == status.ReservationStatusNoShow || reservation.StatusId == status.ReservationStatusPending){
				unattendedCount = unattendedCount + 1
		}
	}
	OneMonth := (time.Hour * 24) * 30
	if(unattendedCount >= MaxAllowedUnattendedReservations){
		latestReservation, err := reservations[0].TransformDateStrToTime()
		if err != nil {
			return true
		}
		latestReservationOneMonthFromNow := latestReservation.Add(OneMonth)
		today := time.Now().Truncate(24 * time.Hour)
		
		if(today.Before(latestReservationOneMonthFromNow)){
			return true
		}
	}
	return false
	
}