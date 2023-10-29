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
	err :=repo.db.Select(&slots,"SELECT id, start_time, end_time, max_capacity FROM time_slots where deleted_at is null")
	return slots, err
}