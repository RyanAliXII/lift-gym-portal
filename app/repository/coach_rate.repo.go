package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)

type CoachRateRepository struct {
	db *sqlx.DB
}

func  NewCoachRateRepository() CoachRateRepository{
	return CoachRateRepository{
		db: db.GetConnection(),
	}
}
func (repo * CoachRateRepository) NewRate(rate model.CoachRate) error {
	_, err := repo.db.Exec("INSERT INTO coaching_rate(description, price, coach_id) VALUES(?, ?, ?)", rate.Description, rate.Price, rate.CoachId)
	return  err
}
func (repo * CoachRateRepository) UpdateRate(rate model.CoachRate) error {
	_, err := repo.db.Exec("UPDATE coaching_rate set description = ?, price = ?, coach_id = ? where id = ? and coach_id = ?", rate.Description, rate.Price, rate.CoachId, rate.Id, rate.CoachId)
	return  err
}
func (repo * CoachRateRepository) DeleteRateByIdAndCoachId(id int, coachId int) error {
	_, err := repo.db.Exec("UPDATE coaching_rate set deleted_at = now() where id = ? and coach_id = ?", id, coachId)
	return  err
}
func (repo * CoachRateRepository) GetRatesByCoachId(id int) ([]model.CoachRate, error){
	rates := make([]model.CoachRate, 0)
	err := repo.db.Select(&rates, "SELECT id, description, price, coach_id from coaching_rate where coach_id = ? and deleted_at is null ORDER BY updated_at DESC", id)
	return rates,err 
}