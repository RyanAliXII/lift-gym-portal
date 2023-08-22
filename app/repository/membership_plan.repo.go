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
	insertQuery := `INSERT INTO membership_plan(description,months, price) VALUES(?, ?, ?)`
	_, insertErr := repo.db.Exec(insertQuery, plan.Description, plan.Months, plan.Price)
	return insertErr
}



func NewMembershipPlanRepository() MembershipPlanRepository{
	return MembershipPlanRepository{
		db: db.GetConnection(),
	}
}
