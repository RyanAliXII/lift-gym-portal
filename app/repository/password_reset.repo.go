package repository

import (
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
	gonanoid "github.com/matoous/go-nanoid/v2"
)


type PasswordReset struct {
	db * sqlx.DB
}

func (repo * PasswordReset) New(accountId int) (model.PasswordReset, error) {
	passwordReset := model.PasswordReset{}
	id, err := gonanoid.New()
	if err != nil {
		return passwordReset, err
	}
	_, err = repo.db.Exec("insert into password_reset(account_id, public_key, expires_at) VALUES(?, ?, DATE_ADD(NOW(), INTERVAL 15 MINUTE))", accountId, id)
	passwordReset.AccountId = accountId
	passwordReset.Id = id
	return passwordReset, err
}

