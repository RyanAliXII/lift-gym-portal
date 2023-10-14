package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)

type ClientLogRepository struct {
	db *sqlx.DB
}
func NewClientLogRepository() ClientLogRepository {
	return ClientLogRepository{
		db: db.GetConnection(),
	}
}
func (repo * ClientLogRepository)NewLog(log model.ClientLog) error{
	_, err := repo.db.Exec("INSERT INTO client_log(client_id, amount_paid, is_member) VALUES(?, ?, ?)", log.ClientId, log.AmountPaid, log.IsMember)
	return err
}