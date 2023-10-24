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
	selectQuery := `SELECT coach.id, coach.given_name, coach.middle_name, coach.surname, 
	coach.date_of_birth, coach.address, coach.emergency_contact,coach.mobile_number, 
	account.email, account.id as account_id, description, COALESCE(CONCAT('[',GROUP_CONCAT('"',coach_image.path,'"'),']'), '[]') as images from coach
	INNER JOIN account on coach.account_id = account.id
	LEFT JOIN coach_image on coach.id = coach_image.coach_id
	GROUP BY coach.id
	ORDER BY coach.updated_at DESC`
	selectErr := repo.db.Select(&coaches, selectQuery)
	return coaches, selectErr 
}
func (repo *CoachRepository)GetCoachById (id int ) (model.Coach, error) {
	coach := model.Coach{}
	selectQuery := `SELECT coach.id, coach.given_name, coach.middle_name, coach.surname, coach.date_of_birth, coach.address, coach.emergency_contact,coach.mobile_number, account.email, account.id as account_id,description from coach
	INNER JOIN account on coach.account_id = account.id where coach.id = ? ORDER BY coach.updated_at DESC LIMIT 1`
	err := repo.db.Get(&coach, selectQuery, id)
	return coach, err
}
func (repo *CoachRepository)UpdateCoachDescription (id int, description string ) (error) {
	updateQuery := `UPDATE coach set description = ? where id = ?`
	_,err := repo.db.Exec( updateQuery, description, id)
	return err
}
func (repo *CoachRepository)GetCoachByIdWithPassword (id int ) (model.Coach, error) {
	coach := model.Coach{}
	selectQuery := `SELECT coach.id, coach.given_name, coach.middle_name, coach.surname, coach.date_of_birth, coach.address, coach.emergency_contact,coach.mobile_number, account.email, account.password, account.id as account_id, description from coach
	INNER JOIN account on coach.account_id = account.id where coach.id = ? ORDER BY coach.updated_at DESC LIMIT 1`
	err := repo.db.Get(&coach, selectQuery, id)
	return coach, err
}
func (repo  CoachRepository)UpdateCoach(coach model.Coach) error {
	dbCoach, getClientErr := repo.GetCoachById(coach.Id)
	if getClientErr != nil {
		return getClientErr
	}
	transaction, transactErr := repo.db.Begin()
	if transactErr != nil {
		transaction.Rollback()
		return transactErr
	}
	updateCoachQuery := `UPDATE coach SET given_name = ?, middle_name = ?, surname = ?, date_of_birth = ?, address = ?, mobile_number = ?, emergency_contact = ? where id = ?`
	_, updateCoachErr := transaction.Exec(updateCoachQuery, coach.GivenName,coach.MiddleName, coach.Surname, coach.DateOfBirth, coach.Address, coach.MobileNumber, coach.EmergencyContact, coach.Id)
	if updateCoachErr != nil {
		transaction.Rollback()
		return updateCoachErr
	}
	updateAccountQuery := `UPDATE account SET email = ? WHERE id = ?`
	_, updateAccountErr := transaction.Exec(updateAccountQuery, coach.Email, dbCoach.AccountId)
	if updateAccountErr != nil {
		transaction.Rollback()
		return updateAccountErr
	}
	transaction.Commit()
	return nil

	
}
func (repo * CoachRepository)UpdatePassword (newPassword string, coachId int )(error){
	coach, err := repo.GetCoachById(coachId)
	if err!= nil {
		return err
	}
	updateQuery := "UPDATE account SET password = ? WHERE id = ?"
	_, err = repo.db.Exec(updateQuery, newPassword, coach.AccountId)
	if err != nil {
		return err
	}
	return nil
} 
func NewCoachRepository()CoachRepository {

	return CoachRepository{
		db: db.GetConnection(),
	}
}