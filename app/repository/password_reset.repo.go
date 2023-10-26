package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
	gonanoid "github.com/matoous/go-nanoid/v2"
)


type PasswordReset struct {
	db * sqlx.DB
}

func NewPasswordResetRepository() PasswordReset{
	return PasswordReset{
		db: db.GetConnection(),
	}
}
func (repo * PasswordReset) New(accountId int) (model.PasswordReset, error) {
	passwordReset := model.PasswordReset{}
	id, err := gonanoid.New()
	if err != nil {
		return passwordReset, err
	}
	transaction, err := repo.db.Beginx()
	if err != nil {
		transaction.Rollback()
		return passwordReset, err
	}
	_, err = transaction.Exec("UPDATE password_reset set completed_at = now() where account_id = ?", accountId)
	if err != nil {
		transaction.Rollback()
		return passwordReset, err
	}
	_, err = transaction.Exec("insert into password_reset(account_id, public_key, expires_at) VALUES(?, ?, DATE_ADD(NOW(), INTERVAL 15 MINUTE))", accountId, id)
	if err != nil {
		transaction.Rollback()
		return passwordReset, err
	}
	passwordReset.AccountId = accountId
	passwordReset.PublicId = id
	
	transaction.Commit()
	return passwordReset, nil
}


func (repo * PasswordReset) GetByPublicKey(publicKey string) (model.PasswordReset, error) {
	passwordReset := model.PasswordReset{}
	err := repo.db.Get(&passwordReset, "SELECT id, public_key, account_id from password_reset where public_key = ? AND expires_at >= NOW() AND completed_at is null",  publicKey )
	return passwordReset, err
}

