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
	repo.db.Select(&slots, "SELECT id, date from date_slot where date >= CAST(CONVERT_TZ(CURDATE(), 'UTC', 'Asia/Manila') as date) and deleted_at is null")
	return slots, nil
}
func (repo * DateSlot) DeleteSlot(id int) (error) {
	_, err := repo.db.Exec("UPDATE date_slot set deleted_at = now() where id = ?", id)
	return err
}