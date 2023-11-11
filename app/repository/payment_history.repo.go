package repository

import (
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)


type PaymentHistory struct {
	db  * sqlx.DB
}
func(repo  * PaymentHistory)GetPaymentHistoryByClient(clientId int )([]model.PaymentHistory, error) {
	payments := make([]model.PaymentHistory, 0)
	query := `SELECT client_id, price as amount, concat('Membership fee: ', membership_plan_snapshot.description) as description, subscription.created_at FROM subscription INNER JOIN membership_plan_snapshot on subscription.membership_plan_snapshot_id = membership_plan_snapshot.id
	UNION ALL
	SELECT client_id, package_snapshot.price as amount, concat('Package payment: ', package_snapshot.description) as description, package_request.created_at FROM package_request INNER JOIN package_snapshot on package_snapshot_id = package_snapshot.id
	UNION ALL
	SELECT  client_id,  coaching_rate_snapshot.price as amount, CONCAT('Coaching fee: ', coaching_rate_snapshot.description) as description, hired_coach.created_at FROM hired_coach INNER JOIN hired_coaches_status on status_id = hired_coaches_status.id  
	INNER JOIN coaching_rate_snapshot on rate_snapshot_id = coaching_rate_snapshot.id
	where client_id = ?
	`
	err := repo.db.Select(&payments, query, clientId)
	return payments,  err
}