package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)


type CoachSchedule struct {
	db * sqlx.DB
}
func NewCoachSchedule()CoachSchedule{
	return CoachSchedule{
		db: db.GetConnection(),
	}
}
func(repo * CoachSchedule) NewSchedule(schedule model.CoachSchedule) error{
	_, err := repo.db.Exec("INSERT INTO coach_schedule(date, time, coach_id) VALUES(?, ?, ?)", schedule.Date, schedule.Time, schedule.CoachId)
	return err
}
func(repo * CoachSchedule) UpdateSchedule(schedule model.CoachSchedule) error{
	_, err := repo.db.Exec("UPDATE coach_schedule set date = ?, time = ? WHERE id = ? and coach_id = ? and deleted_at is null", schedule.Date, schedule.Time, schedule.Id, schedule.CoachId)
	return err
}
func(repo * CoachSchedule) DeleteSchedule(schedule model.CoachSchedule) error{
	_, err := repo.db.Exec("UPDATE coach_schedule set deleted_at = now() WHERE id = ? and coach_id = ? and deleted_at is null", schedule.Id, schedule.CoachId)
	return err
}
func(repo * CoachSchedule)GetSchedulesByCoachId(coachId int ) ([]model.CoachSchedule, error) {
	schedules := make([]model.CoachSchedule, 0)
	err := repo.db.Select(&schedules, `SELECT id, coach_id, time,date from coach_schedule 
	where coach_id = ? and deleted_at is null order by created_at desc`, coachId)
	return schedules, err
}


