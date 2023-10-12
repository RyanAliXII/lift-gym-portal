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
	insertUserQuery := `INSERT INTO user(given_name, middle_name, surname, account_id, role_id) VALUES(?, ?, ?, ?, ?)`
	_, err = transaction.Exec(insertUserQuery, staff.GivenName, staff.MiddleName, staff.Surname, accountId, staff.RoleId)
	if err != nil {
		transaction.Rollback()
		return err
	}
	transaction.Commit()
	return nil	
}

func (repo * StaffRepository) UpdateStaff(staff model.Staff) error{
	transaction, err := repo.db.Beginx()
	if err != nil {
		transaction.Rollback()
		return err
	}
	dbStaff := model.Staff{}
	getQuery := `SELECT account_id FROM user INNER JOIN account on user.account_id = account.id and account.is_root = false where user.id = ? LIMIT 1`
	err = transaction.Get(&dbStaff, getQuery, staff.Id)	
	if err  != nil {
		transaction.Rollback()
		return err
	}
	updateAccountQuery := `UPDATE account set email = ? where id = ?`
	_, err = transaction.Exec(updateAccountQuery, staff.Email, dbStaff.AccountId)
	if err != nil {
		transaction.Rollback()
		return err
	}
	updateUserQuery := `UPDATE user set given_name = ?, middle_name = ?, surname = ? where id = ?`
	_, err = transaction.Exec(updateUserQuery, staff.GivenName, staff.MiddleName, staff.Surname, staff.Id)
	if err != nil {
		transaction.Rollback()
		return err
	}
	transaction.Commit()
	return nil	
}

func (repo *StaffRepository)GetStaffs()([]model.Staff, error){
	staffs := make([]model.Staff, 0)
	query := `SELECT user.id, given_name, middle_name, surname, (case when role_id is null then 0 else role_id end) as role_id,  email FROM user INNER JOIN account on user.account_id = account.id and account.is_root = false ORDER BY user.updated_at DESC`
	err := repo.db.Select(&staffs, query)
	if err != nil {
		 return staffs, err
	}
	return staffs, nil

}
func (repo * StaffRepository) UpdatePassword(newPassword string, userId int) error{
	transaction, err := repo.db.Beginx()
	if err != nil {
		transaction.Rollback()
		return err
	}
	dbStaff := model.Staff{}
	getQuery := `SELECT account_id FROM user INNER JOIN account on user.account_id = account.id and account.is_root = false where user.id = ? LIMIT 1`
	err = transaction.Get(&dbStaff, getQuery, userId)	
	if err  != nil {
		transaction.Rollback()
		return err
	}
	updateAccountQuery := `UPDATE account set password = ? where id = ?`
	_, err = transaction.Exec(updateAccountQuery, newPassword, dbStaff.AccountId)
	if err != nil {
		transaction.Rollback()
		return err
	}
	transaction.Commit()
	return nil	
}
func NewStaffRepository()StaffRepository{
	return StaffRepository{
		db: db.GetConnection(),
	}
}