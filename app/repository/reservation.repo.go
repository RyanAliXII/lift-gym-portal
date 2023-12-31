package repository

import (
	"fmt"
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/status"

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
	err := repo.db.Get(&recordCount,  "SELECT count(1) as recordCount from reservation where client_id = ? and date_slot_id = ? and status_id != 4", reservation.ClientId, reservation.DateSlotId )
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
	 reservation.reservation_id,
	 reservation.id,
	 reservation.client_id, 
	 reservation.date_slot_id, 
	 reservation.time_slot_id,
     (CASE WHEN CURDATE() > date_slot.date AND reservation.status_id = 1 THEN 3 ELSE reservation.status_id END ) AS status_id,
     (CASE WHEN CURDATE() > date_slot.date AND reservation.status_id = 1 THEN "No-Show" ELSE reservation_status.description END) as status,
	 reservation.remarks,
	 JSON_OBJECT('id', client.id, 'givenName', client.given_name, 'middleName', client.middle_name, 'surname', client.surname)  AS client,
	 date_slot.date,
	 CONCAT(TIME_FORMAT(time_slot.start_time, '%h:%i %p'), ' - ', TIME_FORMAT(time_slot.end_time, '%h:%i %p')) AS time
	 FROM reservation 
	 INNER JOIN date_slot ON date_slot_id = date_slot.id
	 INNER JOIN time_slot ON time_slot_id = time_slot.id
	 INNER JOIN client ON reservation.client_id = client.id
	 INNER JOIN reservation_status ON reservation.status_id = reservation_status.id
	 ORDER BY date_slot.date DESC;
	`
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
	 reservation.time_slot_id,
	 reservation.remarks,
	 (CASE WHEN CURDATE() > date_slot.date AND reservation.status_id = 1 THEN 3 ELSE reservation.status_id END ) AS status_id,
     (CASE WHEN CURDATE() > date_slot.date AND reservation.status_id = 1 THEN "No-Show" ELSE reservation_status.description END) as status,
	 JSON_OBJECT('id', client.id, 'givenName', client.given_name, 'middleName', client.middle_name, 'surname', client.surname)  as client,
	 date_slot.date,
	 CONCAT(TIME_format(time_slot.start_time, '%h:%i %p'), ' - ', TIME_format(time_slot.end_time, '%h:%i %p')) as time,
	 reservation_id
	 FROM reservation 
	INNER JOIN date_slot on date_slot_id = date_slot.id
	INNER JOIN time_slot on time_slot_id = time_slot.id
	INNER JOIN client on reservation.client_id = client.id
	INNER JOIN reservation_status on reservation.status_id = reservation_status.id
	where date_slot_id = ? 	ORDER BY date_slot.date DESC
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
	 reservation.time_slot_id,
	 (CASE WHEN CURDATE() > date_slot.date AND reservation.status_id = 1 THEN 3 ELSE reservation.status_id END ) AS status_id,
     (CASE WHEN CURDATE() > date_slot.date AND reservation.status_id = 1 THEN "No-Show" ELSE reservation_status.description END) as status,
	 reservation.remarks,
	 date_slot.date,
	 CONCAT(TIME_format(time_slot.start_time, '%h:%i %p'), ' - ', TIME_format(time_slot.end_time, '%h:%i %p')) as time,
	 reservation_id
	 FROM reservation 
	INNER JOIN date_slot on date_slot_id = date_slot.id
	INNER JOIN time_slot on time_slot_id = time_slot.id
	INNER JOIN reservation_status on reservation.status_id = reservation_status.id
	where client_id = ? ORDER BY reservation.created_at DESC
	`
	err := repo.db.Select(&reservations, query, clientId)
	return reservations, err
}

func (repo * Reservation)MarkAsAttended(id int) error {
	//mark as attended if reservation status has the same id and status is pending or no show.
	_, err := repo.db.Exec("UPDATE reservation set status_id = ? where id = ? and (status_id = 1 or status_id = 3)", status.ReservationStatusAttended, id)
	return err 
}

func (repo * Reservation)MarkAsNoShow(id int) error {
	//mark as attended if reservation status has the same id and status is pending or attended.
	_, err := repo.db.Exec("UPDATE reservation set status_id = ? where id = ? and (status_id = 1 or status_id = 2)", status.ReservationStatusNoShow, id)
	return err 
}

func (repo * Reservation)MarkAsCancelled(id int, remarks string) error {
	//mark as attended if reservation status has the same id and status is pending
	_, err := repo.db.Exec("UPDATE reservation set status_id = ?, remarks = ? where id = ? and status_id = 1", status.ReservationStatusCancelled, remarks, id)
	return err 
}