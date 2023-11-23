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
	query := `SELECT * FROM (SELECT client_id, price as amount, concat('Membership fee: ', membership_plan_snapshot.description) as description, convert_tz(subscription.created_at, 'UTC', 'Asia/Manila') as created_at  FROM subscription INNER JOIN membership_plan_snapshot on subscription.membership_plan_snapshot_id = membership_plan_snapshot.id
	UNION ALL
	SELECT client_id, package_snapshot.price as amount, concat('Package payment: ', package_snapshot.description) as description, convert_tz(package_request.created_at, 'UTC', 'Asia/Manila') as created_at  FROM package_request INNER JOIN package_snapshot on package_snapshot_id = package_snapshot.id
	UNION ALL
	SELECT  client_id,  coaching_rate_snapshot.price as amount, CONCAT('Coaching fee: ', coaching_rate_snapshot.description) as description, convert_tz(hired_coach.created_at, 'UTC', 'Asia/Manila') as created_at FROM hired_coach INNER JOIN hired_coaches_status on status_id = hired_coaches_status.id  
	INNER JOIN coaching_rate_snapshot on rate_snapshot_id = coaching_rate_snapshot.id where status_id = 3
    UNION ALL
    SELECT client_id, amount, CONCAT('Penalty fee: ', 'Not showing up to coach appointment') as description, convert_tz(cap.created_at, 'UTC', 'Asia/Manila') as created_at from coach_appointment_penalty as cap where settled_at is not null
    ) as payments  where payments.client_id = ?
	ORDER BY created_at DESC
	`
	err := repo.db.Select(&payments, query, clientId)
	return payments,  err
}
func(repo  * PaymentHistory)GetPaymentHistory()([]model.PaymentHistory, error) {
	payments := make([]model.PaymentHistory, 0)
	query := `SELECT client_id, amount, description, payment.created_at, JSON_OBJECT('publicId', client.public_id, 'id', client.id, 'givenName', client.given_name, 'middleName', client.middle_name, 'surname', client.surname)  as client FROM 
	(SELECT client_id, price as amount, concat('Membership Payment: ', membership_plan_snapshot.description) as description, convert_tz(subscription.created_at, 'UTC', 'Asia/Manila') as created_at  FROM subscription INNER JOIN membership_plan_snapshot on subscription.membership_plan_snapshot_id = membership_plan_snapshot.id
	 UNION ALL
	 SELECT client_id, package_snapshot.price as amount, concat('Package Payment: ', package_snapshot.description) as description, convert_tz(package_request.created_at, 'UTC', 'Asia/Manila') as created_at  FROM package_request INNER JOIN package_snapshot on package_snapshot_id = package_snapshot.id
	 UNION ALL
	 SELECT  client_id,  coaching_rate_snapshot.price as amount, CONCAT('Coaching payment to ', coach.given_name, ' ', coach.surname,' Coaching rate: ' ,coaching_rate_snapshot.description) as description, convert_tz(hired_coach.created_at, 'UTC', 'Asia/Manila') as created_at FROM hired_coach INNER JOIN hired_coaches_status on status_id = hired_coaches_status.id  
	 INNER JOIN coaching_rate_snapshot on rate_snapshot_id = coaching_rate_snapshot.id
	 INNER JOIN coach ON hired_coach.coach_id = coach.id
	 where status_id = 3
	 ) as payment INNER JOIN client on payment.client_id = client.id
	 ORDER BY created_at DESC
	`
	err := repo.db.Select(&payments, query)
	return payments,  err
}

 
func(repo  * PaymentHistory)GetCoachPayments(coachId int)([]model.PaymentHistory, error) {
	payments := make([]model.PaymentHistory, 0)
	query := `SELECT client_id, description, amount, history.created_at,
	JSON_OBJECT('publicId', client.public_id, 'id', client.id, 'givenName',
	client.given_name, 'middleName', client.middle_name, 'surname', client.surname) AS client
	 FROM (
	  SELECT
		client_id, 
		hired_coach.coach_id,
		coaching_rate_snapshot.price AS amount,
	   'Coaching payment' AS description,
		CONVERT_TZ(hired_coach.created_at, 'UTC', 'Asia/Manila') AS created_at
		FROM
		hired_coach
		INNER JOIN coaching_rate_snapshot ON rate_snapshot_id = coaching_rate_snapshot.id
		INNER JOIN coach ON hired_coach.coach_id = coach.id
		WHERE
		status_id = 3
		UNION ALL
		SELECT 
		client_id, 
		coach_id, 
		amount,
		'Penalty disbursement' as description,
		CONVERT_TZ(coach_appointment_penalty.created_at, 'UTC', 'Asia/Manila') AS created_at 
		from coach_appointment_penalty
	) as history INNER JOIN client ON client_id = client.id
	where history.coach_id = ?
	order by created_at desc
	`
	err := repo.db.Select(&payments, query, coachId)
	return payments,  err
}