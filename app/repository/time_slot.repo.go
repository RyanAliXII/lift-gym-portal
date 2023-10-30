package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)


type TimeSlot struct {
	db * sqlx.DB
}

func NewTimeSlotRepository() TimeSlot {
	return TimeSlot{
		db: db.GetConnection(),
	}
}
func (repo * TimeSlot)NewTimeSlot(timeSlot model.TimeSlot) error {
	_, err :=repo.db.Exec("INSERT INTO time_slot(start_time, end_time, max_capacity) VALUES(?, ?, ?)", timeSlot.StartTime, timeSlot.EndTime, timeSlot.MaxCapacity )
	return err
}
func (repo * TimeSlot)GetTimeSlots() ([]model.TimeSlot, error) {
	slots := make([]model.TimeSlot, 0)
	err :=repo.db.Select(&slots,"SELECT id, start_time, end_time, max_capacity FROM time_slot where deleted_at is null")
	return slots, err
}
func (repo * TimeSlot)DeleteTimeSlot(id int)(error) {
	_, err :=repo.db.Exec("UPDATE time_slot set deleted_at = now() where id = ?", id )
	return err
}

func(repo * TimeSlot) GetTimeSlotsBasedOnDateSlot(dateSlotId int)([]model.TimeSlot, error){
	query := `SELECT time_slot.id, time_slot.start_time, time_slot.end_time, max_capacity, COALESCE(count(reservation.id), 0) as booked, (max_capacity - COALESCE(count(reservation.id), 0)) as available   FROM time_slot 
	LEFT JOIN reservation on time_slot.id = reservation.time_slot_id and reservation.date_slot_id = ?
	where time_slot.deleted_at is null 
	GROUP BY time_slot.id`
	slots := make([]model.TimeSlot, 0)
	err := repo.db.Select(&slots, query, dateSlotId)
	return slots, err
}