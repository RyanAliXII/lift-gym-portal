package repository

import (
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)

type HiredCoachRepository struct {
	db *sqlx.DB
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
		return err
	}
	snapshotId, err = result.LastInsertId()

	return nil
}