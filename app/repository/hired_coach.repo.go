package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)

type HiredCoachRepository struct {
	db *sqlx.DB
}
func NewHiredCoachRepository() (HiredCoachRepository){
	return HiredCoachRepository{
		db: db.GetConnection(),
	}
}
func (repo * HiredCoachRepository) Hire(hiredCoach model.HiredCoach) error {
	transaction, err := repo.db.Beginx()
	if err != nil {
		transaction.Rollback()
		return err
	}
	coachRate := model.CoachRate{}
	err = transaction.Get(&coachRate, "SELECT id, description, price, coach_id from coaching_rate where id = ? and coach_id = ? LIMIT 1", hiredCoach.RateId, hiredCoach.CoachId)
	if err != nil {
		transaction.Rollback()
		return err
	}
	result, err := transaction.Exec("INSERT INTO coaching_rate_snapshot (description, price, coach_id) VALUES(?, ?, ?, ?)", coachRate.Description, coachRate.Price, coachRate.CoachId)
	if err != nil {
		transaction.Rollback()
		return err
	}
	snapshotId, err := result.LastInsertId()
	if err != nil {
		transaction.Rollback()
		return err
	}
	_, err = transaction.Exec("INSERT INTO hired_coach(coach_id, rate_id, rate_snapshot_id, client_id) VALUES(?, ?, ?, ?)", coachRate.CoachId, coachRate.Id, snapshotId, hiredCoach.ClientId)
	if err != nil {
		transaction.Rollback()
		return err
	}
	transaction.Commit()
	return nil
}