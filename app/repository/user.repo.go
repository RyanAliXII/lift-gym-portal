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

func NewUserRepository () UserRepository{
	return UserRepository{
		db: db.GetConnection(),
	}
}