package repository

import (
	"github.com/jmoiron/sqlx"
)


type AccountRepository struct {
	db * sqlx.DB

}

func (repo * AccountRepository ) GetAccountByEmail (email string) error {

	return nil
}