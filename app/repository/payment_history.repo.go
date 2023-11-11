package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)


type PaymentHistory struct {
	db  * sqlx.DB
}

func NewPaymentHistory() PaymentHistory{
	return PaymentHistory{
		db: db.GetConnection(),
	}
}
func(repo  * PaymentHistory)GetPaymentHistoryByClient(clientId int )([]model.PaymentHistory, error) {
	payments := make([]model.PaymentHistory, 0)
	query := `   SELECT * FROM (SELECT client_id, price as amount, concat('Membership fee: ', membership_plan_snapshot.description) as description, convert_tz(subscription.created_at, 'UTC', 'Asia/Manila') as created_at  FROM subscription INNER JOIN membership_plan_snapshot on subscription.membership_plan_snapshot_id = membership_plan_snapshot.id
	UNION ALL
	SELECT client_id, package_snapshot.price as amount, concat('Package payment: ', package_snapshot.description) as description, convert_tz(package_request.created_at, 'UTC', 'Asia/Manila') as created_at  FROM package_request INNER JOIN package_snapshot on package_snapshot_id = package_snapshot.id
	UNION ALL
	SELECT  client_id,  coaching_rate_snapshot.price as amount, CONCAT('Coaching fee: ', coaching_rate_snapshot.description) as description, convert_tz(hired_coach.created_at, 'UTC', 'Asia/Manila') as created_at FROM hired_coach INNER JOIN hired_coaches_status on status_id = hired_coaches_status.id  
	INNER JOIN coaching_rate_snapshot on rate_snapshot_id = coaching_rate_snapshot.id) as payments where payments.client_id = ?

	`
	err := repo.db.Select(&payments, query, clientId)
	return payments,  err
}