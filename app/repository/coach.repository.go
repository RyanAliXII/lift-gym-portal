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
	insertCoachQuery := `INSERT INTO coach(given_name, middle_name, surname, date_of_birth, address, mobile_number, emergency_contact, account_id) VALUES(?, ? ,? , ?, ? , ?, ?, ?)`
	_, insertCoachErr := transaction.Exec(insertCoachQuery, coach.GivenName, coach.MiddleName, coach.Surname, coach.DateOfBirth, coach.Address, coach.MobileNumber, coach.EmergencyContact, accountId )
	if insertCoachErr != nil {
		transaction.Rollback()
		return insertCoachErr
	}
	transaction.Commit()
	return nil
}
func (repo *CoachRepository) GetCoaches() ([]model.Coach, error){
	coaches := make([]model.Coach , 0)
	selectQuery := `SELECT coach.id, coach.given_name, coach.middle_name, coach.surname, coach.date_of_birth, coach.address, coach.emergency_contact,coach.mobile_number, account.email, account.id as account_id from coach
	INNER JOIN account on coach.account_id = account.id ORDER BY coach.updated_at DESC`
	selectErr := repo.db.Select(&coaches, selectQuery)
	return coaches, selectErr 
}
func (repo *CoachRepository)GetCoachById (id int ) (model.Coach, error) {
	coach := model.Coach{}
	selectQuery := `SELECT coach.id, coach.given_name, coach.middle_name, coach.surname, coach.date_of_birth, coach.address, coach.emergency_contact,coach.mobile_number, account.email, account.id as account_id from coach
	INNER JOIN account on coach.account_id = account.id where coach.id = ? ORDER BY coach.updated_at DESC LIMIT 1`
	err := repo.db.Get(&coach, selectQuery, id)
	return coach, err
}

func NewCoachRepository()CoachRepository {

	return CoachRepository{
		db: db.GetConnection(),
	}
}