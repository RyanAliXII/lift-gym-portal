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

func (repo * Reservation)GetReservations () ([]model.Reservation, error){
	reservations := make([]model.Reservation, 0)
	query := `
	SELECT 
	reservation.id,
	 reservation.client_id, 
	 reservation.date_slot_id, 
	 reservation.time_slot_id ,
	 date_slot.date,
	 (case when cancelled_at is null then false else true end) as is_cancelled,
	 (case when attended_at is null then false else true end) as has_attended,
	 CONCAT(TIME_format(time_slot.start_time, '%h:%i %p'), ' - ', TIME_format(time_slot.end_time, '%h:%i %p')) as time,
	 reservation_id
	 FROM reservation 
	INNER JOIN date_slot on date_slot_id = date_slot.id
	INNER JOIN time_slot on time_slot_id = time_slot.id`
	err := repo.db.Select(&reservations, query)
	return reservations, err
}
func (repo * Reservation)GetReservationByDateSlot (dateSlotId int) ([]model.Reservation, error){
	reservations := make([]model.Reservation, 0)
	query := `
	SELECT 
	reservation.id,
	 reservation.client_id, 
	 reservation.date_slot_id, 
	 reservation.time_slot_id ,
	 date_slot.date,
	 (case when cancelled_at is null then false else true end) as is_cancelled,
	 (case when attended_at is null then false else true end) as has_attended,
	 CONCAT(TIME_format(time_slot.start_time, '%h:%i %p'), ' - ', TIME_format(time_slot.end_time, '%h:%i %p')) as time,
	 reservation_id
	 FROM reservation 
	INNER JOIN date_slot on date_slot_id = date_slot.id
	INNER JOIN time_slot on time_slot_id = time_slot.id
	where date_slot_id = ?
	`
	err := repo.db.Select(&reservations, query, dateSlotId)
	return reservations, err
}

func (repo * Reservation)GetClientReservation(clientId int) ([]model.Reservation, error){
	reservations := make([]model.Reservation, 0)
	query := `
	SELECT 
	reservation.id,
	 reservation.client_id, 
	 reservation.date_slot_id, 
	 reservation.time_slot_id ,
	 date_slot.date,
	 (case when cancelled_at is null then false else true end) as is_cancelled,
	 (case when attended_at is null then false else true end) as has_attended,
	 CONCAT(TIME_format(time_slot.start_time, '%h:%i %p'), ' - ', TIME_format(time_slot.end_time, '%h:%i %p')) as time,
	 reservation_id
	 FROM reservation 
	INNER JOIN date_slot on date_slot_id = date_slot.id
	INNER JOIN time_slot on time_slot_id = time_slot.id
	where client_id = ?
	`
	err := repo.db.Select(&reservations, query, clientId)
	return reservations, err
}