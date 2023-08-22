package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)

type MembershipPlanRepository struct {
	db *sqlx.DB
}
func (repo * MembershipPlanRepository) New(plan model.MembershipPlan) error{
	return nil
}



func NewMembershipPlanRepository() MembershipPlanRepository{
	return MembershipPlanRepository{
		db: db.GetConnection(),
	}
}
