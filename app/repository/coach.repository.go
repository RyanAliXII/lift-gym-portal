package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)

type CoachRepository struct {
	db *sqlx.DB
}
func (repo *CoachRepository) NewCoach(coach model.Coach) error{

	transaction, transactErr := repo.db.Begin()
	if transactErr != nil {
		transaction.Rollback()
		return transactErr
	}
	result, insertAccountErr := transaction.Exec("INSERT INTO account (email, password)  VALUES (?, ?)", coach.Email, coach.Password)
	if insertAccountErr != nil {
		transaction.Rollback()
		return insertAccountErr
	}
	accountId, lastInsertedIdErr := result.LastInsertId()
	if lastInsertedIdErr != nil {
		transaction.Rollback()
		return lastInsertedIdErr
	}
	insertClientQuery := `INSERT INTO coach(given_name, middle_name, surname, date_of_birth, address, mobile_number, emergency_contact, account_id) VALUES(?, ? ,? , ?, ? , ?, ?, ?)`
	_, insertClientErr := transaction.Exec(insertClientQuery, coach.GivenName, coach.MiddleName, coach.Surname, coach.DateOfBirth, coach.Address, coach.MobileNumber, coach.EmergencyContact, accountId )
	if insertClientErr != nil {
		transaction.Rollback()
		return insertClientErr
	}
	transaction.Commit()
	return nil
}

func NewCoachRepository()CoachRepository {

	return CoachRepository{
		db: db.GetConnection(),
	}
}