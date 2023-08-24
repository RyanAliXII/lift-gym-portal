package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)




type MemberRepository struct {
	db * sqlx.DB
}
func (repo *MemberRepository)Subscribe (sub model.Subscribe) error {
									

	return nil	
}

func NewMemberRepository() MemberRepository{
	return MemberRepository {
		db: db.GetConnection(),
	}
}