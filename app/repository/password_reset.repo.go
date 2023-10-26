package repository

import (
	"github.com/jmoiron/sqlx"
)


type PasswordReset struct {
		db * sqlx.DB
}

func (repo * PasswordReset) NewPasswordReset() error {


	return nil
}

