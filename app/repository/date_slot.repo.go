package repository

import (
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)


type DateSlot struct {
	db * sqlx.DB
}


func (repo * DateSlot) NewSlots([]model.DateSlot) error {

	return nil
}