package repository

import (
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
	_, err := repo.db.Exec("INSERT INTO reservation(date_slot_id, time_slot_id, client_id) VALUES(?, ?, ?)", reservation.DateSlotId,reservation.TimeSlotId, reservation.ClientId)
	return err
}