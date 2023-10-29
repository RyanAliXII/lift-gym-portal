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
	_, err :=repo.db.Exec("INSERT INTO time_slot(start_time, end_time) VALUES(?, ?)", timeSlot.StartTime, timeSlot.EndTime )
	return err
}