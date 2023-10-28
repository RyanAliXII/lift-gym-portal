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
	_, err := repo.db.NamedExec("INSERT INTO date_slot(date) VALUES(:date) ON DUPLICATE KEY UPDATE date = date", slots)
	return err
}