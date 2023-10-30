package repository

import (
	"fmt"
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)

type Reservation struct {
	db * sqlx.DB
}
func NewReservation()Reservation{
	return Reservation{
		db: db.GetConnection(),
	}
}
func (repo * Reservation)NewReservation(reservation model.Reservation) error {
	recordCount := 1
	err := repo.db.Get(&recordCount,  "SELECT count(1) as recordCount from reservation where client_id = ? and date_slot_id = ?", reservation.ClientId, reservation.DateSlotId )
	if err != nil {
		return err
	}
	if recordCount > 0 {
		return fmt.Errorf("client has an active reservation on this date")
	}
	_, err = repo.db.Exec("INSERT INTO reservation(date_slot_id, time_slot_id, client_id) VALUES(?, ?, ?)", reservation.DateSlotId,reservation.TimeSlotId, reservation.ClientId)
	return err
}