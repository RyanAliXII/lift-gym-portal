package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)

type StaffRepository struct {
	db * sqlx.DB
}

func (repo * StaffRepository) NewStaff(staff model.Staff) error{
	transaction, err := repo.db.Begin()
	if err  != nil {
		transaction.Rollback()
		return err
	}
	insertAccountQuery := `INSERT INTO account(email, password) VALUES(?, ?)`
	result, err := transaction.Exec(insertAccountQuery, staff.Email, staff.Password)
	if err != nil {
		transaction.Rollback()
		return err
	}
	accountId, err:= result.LastInsertId()
	if err  != nil {
		transaction.Rollback()
		return err
	}
	insertUserQuery := `INSERT INTO user(given_name, middle_name, surname, account_id) VALUES(?, ?, ?, ?)`
	_, err = transaction.Exec(insertUserQuery, staff.GivenName, staff.MiddleName, staff.Surname, accountId)
	if err != nil {
		transaction.Rollback()
		return err
	}
	transaction.Commit()
	return nil	
}

func (repo *StaffRepository)GetStaffs()([]model.Staff, error){
	staffs := make([]model.Staff, 0)
	query := `SELECT user.id, given_name, middle_name, surname, email, password FROM user INNER JOIN account on user.account_id = account.id and account.is_root = false`
	err := repo.db.Select(&staffs, query)
	if err != nil {
		 return staffs, err
	}
	return staffs, nil

}
func NewStaffRepository()StaffRepository{
	return StaffRepository{
		db: db.GetConnection(),
	}
}