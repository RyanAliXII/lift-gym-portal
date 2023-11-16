package repository

import (
	"fmt"
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
	insertUserQuery := `INSERT INTO user(given_name, middle_name, surname, account_id, role_id, date_of_birth, gender, address,mobile_number,emergency_contact) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = transaction.Exec(insertUserQuery, staff.GivenName, staff.MiddleName, staff.Surname, accountId, staff.RoleId, staff.DateOfBirth,staff.Gender, staff.Address, staff.MobileNumber, staff.EmergencyContact)
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
	updateUserQuery := `UPDATE user set given_name = ?, middle_name = ?, surname = ?, role_id = ?, date_of_birth = ?, gender = ?, address = ?, mobile_number = ?, emergency_contact = ? where id = ?`
	_, err = transaction.Exec(updateUserQuery, staff.GivenName, staff.MiddleName, staff.Surname, staff.RoleId, staff.DateOfBirth, staff.Gender, staff.Address, staff.MobileNumber, staff.EmergencyContact, staff.Id)
if err != nil {
		transaction.Rollback()
		return err
	}
	transaction.Commit()
	return nil	
}

func (repo * StaffRepository) Delete (id int) error {
	query := `SELECT user.id, given_name, middle_name, surname, (case when role_id is null then 0 else role_id end) as role_id,  email, COALESCE(date_of_birth, '') as date_of_birth, gender, address,mobile_number,emergency_contact, public_id FROM user INNER JOIN account on user.account_id = account.id and account.is_root = false where user.deleted_at is null and user.id = ?  ORDER BY user.updated_at DESC LIMIT 1`
	staff := model.Staff{}
	err := repo.db.Get(&staff, query, id )
	if err != nil {
		return err
	}
	_, err = repo.db.Exec("UPDATE user set deleted_at = now() where id = ?", id)
	return err
}

func (repo *StaffRepository)GetStaffs()([]model.Staff, error){
	staffs := make([]model.Staff, 0)
	query := `SELECT user.id, given_name, middle_name, surname, (case when role_id is null then 0 else role_id end) as role_id,  email, COALESCE(date_of_birth, '') as date_of_birth, gender, address,mobile_number,emergency_contact, public_id FROM user INNER JOIN account on user.account_id = account.id and account.is_root = false where user.deleted_at is null ORDER BY user.updated_at DESC`
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

func (repo * StaffRepository)Search(keyword string)([]model.Staff, error) {
	staffs:= make([]model.Staff , 0)
	keywordLike := fmt.Sprintf("%s%s%s","%", keyword, "%")
	selectQuery := `SELECT staff.id, staff.given_name, staff.middle_name, staff.surname,
	staff.date_of_birth, staff.address, staff.emergency_contact,staff.mobile_number, staff.gender, staff.public_id,
	account.email, account.id as account_id from user as staff
	INNER JOIN account on staff.account_id = account.id
	where (staff.given_name LIKE ? OR staff.middle_name LIKE ? OR staff.surname LIKE ? OR staff.mobile_number LIKE ? OR account.email LIKE ? OR staff.public_id LIKE ?) and staff.deleted_at is null and is_root = false
	ORDER BY staff.updated_at DESC LIMIT 50`
	selectErr := repo.db.Select(&staffs, selectQuery, keywordLike, keywordLike, keywordLike, keywordLike, keywordLike, keywordLike)
	return staffs, selectErr 
}
func NewStaffRepository()StaffRepository{
	return StaffRepository{
		db: db.GetConnection(),
	}
}