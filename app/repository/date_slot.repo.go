package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)


type DateSlot struct {
	db * sqlx.DB
}
func NewDateSlotRepository () DateSlot{
	return DateSlot{
		db: db.GetConnection(),
	}
}
func (repo * DateSlot) NewSlots(slots []model.DateSlot) error {
	_, err := repo.db.NamedExec("INSERT INTO date_slot(date) VALUES(:date) ON DUPLICATE KEY UPDATE deleted_at = null", slots)
	return err
}
func (repo * DateSlot) GetSlots() ( []model.DateSlot, error) {
    slots :=   make([]model.DateSlot, 0)
	repo.db.Select(&slots, `SELECT date_slot.id, date, COUNT(reservation.id) as booked,
	COALESCE(time_slot.total_capacity, 0) as total_capacity,
	COALESCE(time_slot.total_capacity, 0) - COUNT(reservation.id) as available
	 from date_slot
	LEFT JOIN (
		SELECT reservation.id, reservation.date_slot_id, reservation.status_id FROM reservation 
		INNER JOIN time_slot on reservation.time_slot_id = time_slot.id and time_slot.deleted_at is null
    ) as  reservation on date_slot.id = reservation.date_slot_id and reservation.status_id != 4
	INNER JOIN (SELECT SUM(max_capacity) as total_capacity  FROM time_slot where deleted_at is null) as time_slot on true = true
	where date >= CAST(CONVERT_TZ(CURDATE(), 'UTC', 'Asia/Manila') as date) and date_slot.deleted_at is null
	GROUP BY date_slot.id, time_slot.total_capacity`)
	return slots, nil
}

func (repo * DateSlot) GetAllSlots() ( []model.DateSlot, error) {
    slots :=   make([]model.DateSlot, 0)
	repo.db.Select(&slots, `SELECT date_slot.id, date, COUNT(reservation.id) as booked,
	COALESCE(time_slot.total_capacity, 0) as total_capacity,
	COALESCE(time_slot.total_capacity, 0) - COUNT(reservation.id) as available
	 from date_slot
	LEFT JOIN (
		SELECT reservation.id, reservation.date_slot_id, reservation.status_id FROM reservation 
		INNER JOIN time_slot on reservation.time_slot_id = time_slot.id and time_slot.deleted_at is null
    ) as  reservation on date_slot.id = reservation.date_slot_id and reservation.status_id != 4
	INNER JOIN (SELECT SUM(max_capacity) as total_capacity  FROM time_slot where deleted_at is null) as time_slot on true = true
	where date_slot.deleted_at is null
	GROUP BY date_slot.id, time_slot.total_capacity`)
	return slots, nil
}
func (repo * DateSlot) DeleteSlot(id int) (error) {
	_, err := repo.db.Exec("UPDATE date_slot set deleted_at = now() where id = ?", id)
	return err
}
