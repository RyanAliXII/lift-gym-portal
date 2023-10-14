package repository

import (
	"lift-fitness-gym/app/db"

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
func (repo * ClientLogRepository)NewLog(){
	
}