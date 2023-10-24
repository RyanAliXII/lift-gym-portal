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
	SELECT client.id,subscription.id as subscription_id, client.given_name,client.middle_name, client.surname, account.email, client.mobile_number, subscription.valid_until, 
	JSON_OBJECT('id', membership_plan.id, 'description', membership_plan.description, 'months', membership_plan.months, 'price', membership_plan.price) as membership_plan, 
	JSON_OBJECT('id', membership_plan_snapshot.id, 'description', membership_plan_snapshot.description, 'months', membership_plan_snapshot.months, 'price', membership_plan_snapshot.price) as membership_plan_snaphot,
	subscription.created_at FROM subscription
	INNER JOIN client on subscription.client_id = client.id
	INNER JOIN account on client.account_id = account.id
	INNER JOIN membership_plan on subscription.membership_plan_id = membership_plan.id
	INNER JOIN membership_plan_snapshot on subscription.membership_plan_snapshot_id = membership_plan_snapshot.id
	where subscription.valid_until >= NOW() and subscription.cancelled_at is NULL
	ORDER BY subscription.created_at DESC;
	`
	selectErr := repo.db.Select(&members, selectQuery)
	return members, selectErr
}
func (repo * MemberRepository)GetMemberById(id int)(model.Member, error){
	member  := model.Member{}
	query := `SELECT client.id,subscription.id as subscription_id, client.given_name,client.middle_name, client.surname, account.email, client.mobile_number, subscription.valid_until, 
	JSON_OBJECT('id', membership_plan.id, 'description', membership_plan.description, 'months', membership_plan.months, 'price', membership_plan.price) as membership_plan, 
	JSON_OBJECT('id', membership_plan_snapshot.id, 'description', membership_plan_snapshot.description, 'months', membership_plan_snapshot.months, 'price', membership_plan_snapshot.price) as membership_plan_snaphot,
	subscription.created_at FROM subscription
	INNER JOIN client on subscription.client_id = client.id
	INNER JOIN account on client.account_id = account.id
	INNER JOIN membership_plan on subscription.membership_plan_id = membership_plan.id
	INNER JOIN membership_plan_snapshot on subscription.membership_plan_snapshot_id = membership_plan_snapshot.id
	where subscription.valid_until >= NOW() and subscription.cancelled_at is NULL and client.id = ?
	ORDER BY subscription.created_at DESC
	`
	err := repo.db.Get(&member, query, id)
	return member, err
}
func (repo *MemberRepository)Subscribe(sub model.Subscribe) error {
	transaction, transactErr := repo.db.Beginx()
	if transactErr != nil {
		transaction.Rollback()
		return transactErr
	}
	plan := model.MembershipPlan{}
	recordCount := 0
	//check if client has an active subscription
	checkActiveSubErr := transaction.Get(&recordCount, "SELECT COUNT(1) as recordCount FROM subscription WHERE client_id = ? and subscription.valid_until >= NOW() and cancelled_at is null", sub.ClientId)
	if checkActiveSubErr != nil {
		transaction.Rollback()
		return checkActiveSubErr
	}
	if recordCount > 0 {
		transaction.Rollback()
		return fmt.Errorf("client has an active subscription.")
	}
	getErr := transaction.Get(&plan, "SELECT price, description, months from membership_plan where id = ?", sub.MembershipPlanId)
	if getErr != nil {
		transaction.Rollback()
		return getErr
	}
	if (sub.MembershipSnapshotId == 0) {
		result , err := transaction.Exec("INSERT INTO membership_plan_snapshot(description,months, price) VALUES(?, ?, ?)", plan.Description, plan.Months, plan.Price)
		if err != nil {
			return err
		}
		snapshotId, err := result.LastInsertId()
		if err != nil {
			return err
		}
		sub.MembershipSnapshotId = int(snapshotId)
	}
	insertQuery := `INSERT INTO subscription (client_id, membership_plan_id, valid_until, membership_plan_snapshot_id)
	VALUES (?, ?, DATE_ADD(CAST(NOW() AS DATE), INTERVAL ? MONTH), ?)`
	_, insertErr := transaction.Exec(insertQuery, sub.ClientId, sub.MembershipPlanId, plan.Months, sub.MembershipSnapshotId)

	if insertErr != nil {
		transaction.Rollback()
		return insertErr
	}
	transaction.Commit()
	return insertErr
}
func (repo * MemberRepository)CancelSubscription(subId int) error {
	updateQuery := `UPDATE subscription SET cancelled_at = NOW() where id = ?`
	_, updateErr := repo.db.Exec(updateQuery, subId)
	fmt.Println("CANCELLED")
	return updateErr
}
func NewMemberRepository() MemberRepository{
	return MemberRepository {
		db: db.GetConnection(),
	}
}