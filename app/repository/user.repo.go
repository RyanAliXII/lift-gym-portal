package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func (repo *UserRepository) GetUserByEmail(email string)(model.User, error) {
	user := model.User{}
	selectQuery := `SELECT user.id, given_name, middle_name, surname, email, password FROM user INNER JOIN account on user.account_id = account.id where UPPER(email) = UPPER(?) LIMIT 1`
	getErr := repo.db.Get(&user, selectQuery, email)
	return user, getErr

}
func (repo *UserRepository) GetClientUserByEmail (email string) (model.User, error){
	user := model.User{}
	selectQuery := `SELECT client.id, client.given_name,client.middle_name,client.surname, account.email, account.password FROM client 
	INNER JOIN account on client.account_id = account.id where UPPER(account.email) = UPPER(?) LIMIT 1;
	`
	getErr := repo.db.Get(&user, selectQuery, email)
	return user, getErr
}
func (repo  * UserRepository) GetCoachUserByEmail(email string)(model.Coach, error){
	coach := model.Coach{}
	query := `SELECT coach.id, coach.given_name, coach.middle_name, coach.surname, coach.date_of_birth, coach.address, coach.emergency_contact,coach.mobile_number, account.email, account.id as account_id, account.password from coach
	INNER JOIN account on coach.account_id = account.id where UPPER(account.email) = UPPER(?) ORDER BY coach.updated_at DESC`
	err := repo.db.Get(&coach, query, email)
	return coach, err
}
func NewUserRepository () UserRepository{
	return UserRepository{
		db: db.GetConnection(),
	}
}