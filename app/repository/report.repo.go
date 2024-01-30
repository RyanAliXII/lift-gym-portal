package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)


type Report struct {
	db * sqlx.DB
}
func NewReport() Report {
	return Report {
		db: db.GetConnection(),
	}
}
func(repo * Report) GenerateReportData(startDate string, endDate string, preparedBy string)(model.ReportData, error) {
	query := `
	SELECT 
	(SELECT count(1) from client where deleted_at is null and date(created_at) between ? and ?) as clients,
	(SELECT COUNT(1) 
	FROM client 
	INNER JOIN subscription ON subscription.client_id = client.id  
	WHERE subscription.valid_until >= NOW() AND subscription.cancelled_at IS NULL and client.deleted_at is null and  date(subscription.created_at) between ? and ? ) AS members,
	(SELECT COUNT(1) from coach where deleted_at is null and date(created_at) between ? and ? ) as coaches,
	(SELECT  COUNT(1) from package_request where date(created_at) between ? and ? ) as package_requests,
	(SELECT COUNT(1) from membership_request where date(created_at) between ? and ?  ) as membership_requests,
	(SELECT COUNT(1) from reservation where date(created_at) between ? and ?  ) as reservations,
	(SELECT COUNT(1) from equipment where date(created_at) between ? and ?  ) as inventory_items,
	(
		COALESCE(
			(SELECT SUM(amount_paid) 
			FROM client_log  
			WHERE deleted_at  is null and DATE(created_at) BETWEEN ? AND ?), 0) 
		+
		COALESCE(
			(SELECT SUM(price) 
			FROM subscription
			INNER JOIN membership_plan_snapshot 
			ON subscription.membership_plan_snapshot_id = membership_plan_snapshot.id	
			WHERE subscription.valid_until >= NOW() 
			AND subscription.cancelled_at IS NULL 
			AND DATE(subscription.created_at) BETWEEN ? AND ?), 0)
		+
		COALESCE(
			(SELECT SUM(price) 
			FROM package_request 
			INNER JOIN package_snapshot 
			ON package_request.package_snapshot_id = package_snapshot.id 
			WHERE (DATE(package_request.created_at) BETWEEN ? AND ? ) and status_id = 3), 0)
	) AS earnings,
	JSON_OBJECT(
		'walkIn', COALESCE((SELECT SUM(amount_paid) 
		FROM client_log  
		WHERE deleted_at is null and date(created_at) BETWEEN ? AND ?), 0), 
		
		'membership',  COALESCE((SELECT  SUM(price) from subscription
		INNER JOIN membership_plan_snapshot on subscription.membership_plan_snapshot_id = membership_plan_snapshot.id
		where subscription.valid_until >= NOW() 
		and subscription.cancelled_at is NULL and date(subscription.created_at)  BETWEEN ? AND ?),0),
		'package', COALESCE((
			SELECT  SUM(package_snapshot.price) from package_request 
			INNER JOIN package_snapshot on package_snapshot_id = package_snapshot.id
			where date(package_request.created_at) BETWEEN ?
			AND ? and status_id = 3 ),0)
	) as earnings_breakdown;
	`
	data := model.ReportData{}
	data.PreparedBy = preparedBy
	data.StartDate = startDate
	data.EndDate = endDate
	err := repo.db.Get(&data, query,
		startDate, endDate,
		startDate, endDate, 
		startDate, endDate,
		startDate, endDate, 
		startDate, endDate, 
		startDate,endDate, 
		startDate, endDate, 
		startDate, endDate, 
		startDate,endDate,
		startDate,endDate,
		startDate, endDate, 
		startDate, endDate, 
		startDate,endDate,
	)
	if err != nil {
		return model.ReportData{}, nil
	}
	query = `SELECT 
	client_log.id, 
	client_log.client_id, 
	client_log.amount_paid, 
	client_log.is_member, 
	JSON_OBJECT('publicId',client.public_id ,'id', client.id, 'givenName', client.given_name, 'middleName', client.middle_name, 'surname', client.surname, 'email', account.email)  as client,  
	convert_tz(client_log.created_at, 'UTC', 'Asia/Manila') as created_at
	FROM client_log 
	INNER JOIN client ON client_log.client_id = client.id 
	INNER JOIN account ON client.account_id = account.id
	WHERE client_log.created_at BETWEEN ? AND ?
	ORDER BY created_at DESC
	`
	err = repo.db.Select(&data.ClientLogs, query, startDate, endDate)
	if err != nil {
		return model.ReportData{}, err
	}
	query = `
	SELECT client.id,subscription.id as subscription_id, client.given_name,client.middle_name, 
	client.surname, 
	account.email, client.mobile_number, 
	client.public_id, subscription.valid_until, 
	JSON_OBJECT('id', membership_plan_snapshot.id, 'description',
    membership_plan_snapshot.description, 'months', 
    membership_plan_snapshot.months, 'price', membership_plan_snapshot.price) as membership_plan,
	subscription.created_at FROM subscription
	INNER JOIN client on subscription.client_id = client.id
	INNER JOIN account on client.account_id = account.id
	INNER JOIN membership_plan_snapshot on subscription.membership_plan_snapshot_id = membership_plan_snapshot.id
	where subscription.valid_until >= NOW() and subscription.cancelled_at is NULL and subscription.created_at between ? and ?
	ORDER BY subscription.created_at DESC
	`
	err = repo.db.Select(&data.NewMembers, query, startDate, endDate)
	if err != nil {
		return model.ReportData{}, err
	}
	query = `
	SELECT
	coach_appointment_penalty.id, 
	client_id, 
	coach_id, 
	amount,
	(case when settled_at is null then false else true end) as is_settled,  
	JSON_OBJECT(
	   'publicId', client.public_id, 
	   'id', client.id, 
	   'givenName', client.given_name, 
	   'middleName', client.middle_name, 
	   'surname', client.surname)  as client,
	JSON_OBJECT(
	   'id', coach.id,
	   'givenName', coach.given_name,
	   'middleName', coach.middle_name,
	   'surname', coach.surname
	) as coach
	FROM coach_appointment_penalty
	INNER JOIN client on coach_appointment_penalty.client_id = client.id
	INNER JOIN coach on coach_appointment_penalty.coach_id = coach.id
	where coach_appointment_penalty.created_at between ? and ?
	order by coach_appointment_penalty.created_at desc
	`
	err = repo.db.Select(&data.CoachingPenalties, query, startDate, endDate)
	if err != nil {
		return model.ReportData{}, err
	}
	result, err := repo.db.Exec("INSERT INTO report(data) VALUES(?)", data)

	if err != nil {
		return model.ReportData{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return model.ReportData{}, err
	}
	data.Id = id
	return data, nil
}

func (repo * Report) GetReportById (id int) (model.ReportDataJSON, error ){
	data := model.ReportDataJSON{}
	err := repo.db.Get(&data, "Select data from report where id = ? limit 1", id)
	return data, err
}
func (repo *Report)GetWalkIns(startDate string, endDate string)([]model.WalkInData, error) {
	data := make([]model.WalkInData, 0)
	query := `SELECT COUNT(1) as total, cast(convert_tz(created_at, 'UTC', 'Asia/Manila') as date)  as date FROM client_log where date(created_at) BETWEEN ? AND ? GROUP BY DAY(created_at), cast(convert_tz(created_at, 'UTC', 'Asia/Manila') as date) `
	err := repo.db.Select(&data, query, startDate, endDate)
	return data, err
}