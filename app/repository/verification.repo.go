package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type VerificationRepository struct {
	db *sqlx.DB
}

func (repo *VerificationRepository) CreateEmailVerification(clientId string) (model.EmailVerification, error) {
	verification := model.EmailVerification{}
	id, newIdErr := gonanoid.New()
	if newIdErr != nil {
		return verification, newIdErr
	}
	transaction, transactErr := repo.db.Beginx()
	if transactErr != nil {
		transaction.Rollback()
		return verification, transactErr
	}
	insertQuery := "INSERT INTO email_verification(public_id, client_id, expires_at)VALUES(?, ?, DATE_ADD(NOW(), INTERVAL 5 MINUTE))"
	_, insertErr := transaction.Exec(insertQuery, id, clientId)
	if insertErr != nil {
		transaction.Rollback()
		return verification,insertErr
	}
	selectQuery := "Select public_id, client_id, expires_at, created_at  from email_verification where public_id = ?"
	getErr := transaction.Get(&verification, selectQuery, id)
	if getErr != nil {
		transaction.Rollback()
		return verification, getErr
	}

	transaction.Commit()
	return verification, nil
}


func NewVerificationRepository()VerificationRepository{

	return VerificationRepository{
		db: db.GetConnection(),
	}
}