package repository

import (
	"fmt"
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)




type MemberRepository struct {
	db * sqlx.DB
}

func (repo * MemberRepository)GetMembers()([]model.Member,  error){
	members := make([]model.Member, 0)
	selectQuery := `
	SELECT client.id, client.given_name,client.middle_name, client.surname, account.email, client.mobile_number, subscription.valid_until, JSON_OBJECT('id', membership_plan.id, 'description', membership_plan.description, 'months', membership_plan.months, 'price', membership_plan.price) as membership_plan, 
	subscription.created_at FROM subscription
	INNER JOIN client on subscription.client_id = client.id
	INNER JOIN account on client.account_id = account.id
	INNER JOIN membership_plan on subscription.membership_plan_id = membership_plan.id
	where subscription.valid_until >= NOW()
	ORDER BY subscription.created_at DESC
	`
	selectErr := repo.db.Select(&members, selectQuery)
	return members, selectErr
}
func (repo *MemberRepository)Subscribe (sub model.Subscribe) error {
	transaction, transactErr := repo.db.Beginx()
	if transactErr != nil {
		transaction.Rollback()
		return transactErr
	}
	plan := model.MembershipPlan{}
	recordCount := 0
	//check if has an active subscription
	checkActiveSubErr := transaction.Get(&recordCount, "SELECT COUNT(1) as recordCount FROM subscription WHERE client_id = ? and subscription.valid_until >= NOW()", sub.ClientId)
	if checkActiveSubErr != nil {
		transaction.Rollback()
		return checkActiveSubErr
	}
	if recordCount > 0 {
		return fmt.Errorf("client has an active subscription.")
	}
	getErr := transaction.Get(&plan, "SELECT months from membership_plan where id = ?", sub.MembershipPlanId)
	if getErr != nil {
		transaction.Rollback()
		return getErr
	}
	insertQuery := `INSERT INTO subscription (client_id, membership_plan_id, valid_until)
	VALUES (?, ?, DATE_ADD(CAST(NOW() AS DATE), INTERVAL ? MONTH))`
	_, insertErr := transaction.Exec(insertQuery, sub.ClientId, sub.MembershipPlanId, plan.Months)

	if insertErr != nil {
		transaction.Rollback()
		return insertErr
	}
	transaction.Commit()
	return insertErr
}

func NewMemberRepository() MemberRepository{
	return MemberRepository {
		db: db.GetConnection(),
	}
}