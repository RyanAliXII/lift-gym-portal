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
func (repo * DateSlot) NewSlots([]model.DateSlot) error {
	return nil
}