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
			WHERE DATE(package_request.created_at) BETWEEN ? AND ?), 0)
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