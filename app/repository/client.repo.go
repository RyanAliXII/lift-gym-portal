package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)

type ClientRepository struct {
	db *sqlx.DB
}

func (repo *ClientRepository) New(client model.Client)(error) {
		transaction, transactErr := repo.db.Begin()
		if transactErr != nil {
			transaction.Rollback()
			return transactErr
		}
		result, insertAccountErr := transaction.Exec("INSERT INTO account (email, password)  VALUES (?, ?)", client.Email, client.Password)
		if insertAccountErr != nil {
			transaction.Rollback()
			return insertAccountErr
		}
		accountId, lastInsertedIdErr := result.LastInsertId()
		if lastInsertedIdErr != nil {
			transaction.Rollback()
			return lastInsertedIdErr
		}
		insertClientQuery := `INSERT INTO client(given_name, middle_name, surname, date_of_birth, address, mobile_number, emergency_contact, account_id) VALUES(?, ? ,? , ?, ? , ?, ?, ?)`
		_, insertClientErr := transaction.Exec(insertClientQuery, client.GivenName, client.MiddleName, client.Surname, client.DateOfBirth, client.Address, client.MobileNumber, client.EmergencyContact, accountId )
		if insertClientErr != nil {
			transaction.Rollback()
			return insertClientErr
		}
		transaction.Commit()
		return nil
}

func NewClientRepository () ClientRepository{
	return ClientRepository{
		db: db.GetConnection(),
	}
}
