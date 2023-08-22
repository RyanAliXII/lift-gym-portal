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

func (repo  MembershipPlanRepository)Get() ([]model.MembershipPlan, error) {
	plans := make([]model.MembershipPlan, 0)
	selectErr := repo.db.Select(&plans, "SELECT id, description, price, months from membership_plan order by updated_at desc")
	return plans, selectErr

}
func (repo  MembershipPlanRepository)Update(plan model.MembershipPlan)(error) {
	_, updateErr := repo.db.Exec(`UPDATE membership_plan SET description = ?, price = ?, months = ? WHERE id = ?`, plan.Description, plan.Price, plan.Months, plan.Id )
	return updateErr
}



func NewMembershipPlanRepository() MembershipPlanRepository{
	return MembershipPlanRepository{
		db: db.GetConnection(),
	}
}
