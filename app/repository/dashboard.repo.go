package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)


type Dashboard struct {
	db * sqlx.DB
}

func NewDashboard() Dashboard{
	return Dashboard{
		db: db.GetConnection(),
	}
}

func (repo * Dashboard)GetAdminDashboardData() (model.AdminDashboardData, error) {
	data := model.AdminDashboardData{}
	query := `
	SELECT  
    (SELECT COUNT(1) FROM client where deleted_at is null ) AS clients, 
    (SELECT COUNT(1) 
     FROM client 
     INNER JOIN subscription ON subscription.client_id = client.id  
     WHERE subscription.valid_until >= NOW() AND subscription.cancelled_at IS NULL and client.deleted_at is null) AS members,
	
     (
     COALESCE((SELECT SUM(amount_paid) 
     FROM client_log  
     WHERE deleted_at is null and created_at BETWEEN DATE_SUB(NOW(), INTERVAL 1 YEAR) AND NOW()), 0) 
     + 
     COALESCE((SELECT  SUM(price) from subscription
     INNER JOIN membership_plan_snapshot on subscription.membership_plan_snapshot_id = membership_plan_snapshot.id
	  where subscription.valid_until >= NOW() 
	  and subscription.cancelled_at is NULL and subscription.created_at BETWEEN DATE_SUB(NOW(), INTERVAL 1 YEAR) 
      AND NOW()),0)
      + 
      COALESCE((
      SELECT SUM(price) from package_request 
      INNER JOIN package_snapshot on package_request.package_snapshot_id = package_snapshot_id 
      where package_request.created_at BETWEEN DATE_SUB(NOW(), INTERVAL 1 YEAR) 
      AND NOW()),0)
	  ) as annual_earnings,
    (
     COALESCE((SELECT SUM(amount_paid) 
     FROM client_log  
     WHERE created_at BETWEEN DATE_SUB(NOW(), INTERVAL 1 MONTH) AND NOW()), 0) 
     + 
     COALESCE((SELECT  SUM(price) from subscription
     INNER JOIN membership_plan_snapshot on subscription.membership_plan_snapshot_id = membership_plan_snapshot.id
	  where subscription.valid_until >= NOW() 
	  and subscription.cancelled_at is NULL and subscription.created_at BETWEEN DATE_SUB(NOW(), INTERVAL 1 MONTH) 
      AND NOW()),0)
      + 
      COALESCE((
      SELECT SUM(price) from package_request 
      INNER JOIN package_snapshot on package_request.package_snapshot_id = package_snapshot_id 
      where package_request.created_at BETWEEN DATE_SUB(NOW(), INTERVAL 1 MONTH) 
      AND NOW()),0)
	  ) as monthly_earnings,
      
	 (
     COALESCE((SELECT SUM(amount_paid) 
     FROM client_log  
     WHERE deleted_at is null and created_at BETWEEN DATE_SUB(NOW(), INTERVAL 7 DAY) AND NOW()), 0) 
     + 
     COALESCE((SELECT  SUM(price) from subscription
     INNER JOIN membership_plan_snapshot on subscription.membership_plan_snapshot_id = membership_plan_snapshot.id
	  where subscription.valid_until >= NOW() 
	  and subscription.cancelled_at is NULL and subscription.created_at BETWEEN DATE_SUB(NOW(), INTERVAL 7 DAY) 
      AND NOW()),0)
      + 
      COALESCE((
      SELECT SUM(price) from package_request 
      INNER JOIN package_snapshot on package_request.package_snapshot_id = package_snapshot_id 
      where package_request.created_at BETWEEN DATE_SUB(NOW(), INTERVAL 7 DAY) 
      AND NOW()),0)
	  ) as weekly_earnings,
      
	  JSON_OBJECT(
      'walkIn', COALESCE((SELECT SUM(amount_paid) 
      FROM client_log  
      WHERE deleted_at is null and created_at BETWEEN DATE_SUB(NOW(), INTERVAL 1 YEAR) AND NOW()), 0), 
	 'membership',  COALESCE((SELECT  SUM(price) from subscription
      INNER JOIN membership_plan_snapshot on subscription.membership_plan_snapshot_id = membership_plan_snapshot.id
	  where subscription.valid_until >= NOW() 
	  and subscription.cancelled_at is NULL and subscription.created_at BETWEEN DATE_SUB(NOW(), INTERVAL 1 YEAR) 
      AND NOW()),0),
      'package', COALESCE((SELECT SUM(price) from package_request 
      INNER JOIN package_snapshot on package_request.package_snapshot_id = package_snapshot_id 
      where package_request.created_at BETWEEN DATE_SUB(NOW(), INTERVAL 1 YEAR) 
      AND NOW()),0)
	  ) as annual_earnings_breakdown,
      
      JSON_OBJECT(
      'walkIn', COALESCE((SELECT SUM(amount_paid) 
      FROM client_log  
      WHERE deleted_at is null and created_at BETWEEN DATE_SUB(NOW(), INTERVAL 1 MONTH) AND NOW()), 0), 
	 'membership',  COALESCE((SELECT  SUM(price) from subscription
      INNER JOIN membership_plan_snapshot on subscription.membership_plan_snapshot_id = membership_plan_snapshot.id
	  where subscription.valid_until >= NOW() 
	  and subscription.cancelled_at is NULL and subscription.created_at BETWEEN DATE_SUB(NOW(), INTERVAL 1 MONTH) 
      AND NOW()),0),
      'package', COALESCE((SELECT SUM(price) from package_request 
      INNER JOIN package_snapshot on package_request.package_snapshot_id = package_snapshot_id 
      where package_request.created_at BETWEEN DATE_SUB(NOW(), INTERVAL 1 MONTH) 
      AND NOW()),0)
	  ) as monthly_earnings_breakdown,
      
      JSON_OBJECT(	
      'walkIn', COALESCE((SELECT SUM(amount_paid) 
      FROM client_log  
      WHERE created_at BETWEEN DATE_SUB(NOW(), INTERVAL 7 DAY) AND NOW()), 0), 
	 'membership',  COALESCE((SELECT  SUM(price) from subscription
      INNER JOIN membership_plan_snapshot on subscription.membership_plan_snapshot_id = membership_plan_snapshot.id
	  where subscription.valid_until >= NOW() 
	  and subscription.cancelled_at is NULL and subscription.created_at BETWEEN DATE_SUB(NOW(), INTERVAL 7 DAY) 
      AND NOW()),0),
      'package', COALESCE((SELECT SUM(price) from package_request 
      INNER JOIN package_snapshot on package_request.package_snapshot_id = package_snapshot_id 
      where package_request.created_at BETWEEN DATE_SUB(NOW(), INTERVAL 7 DAY) 
      AND NOW()),0)
	  ) as weekly_earnings_breakdown;	   
	`
	err := repo.db.Get(&data, query)
	return data, err
}

func (repo *Dashboard)GetWeeklyWalkIns()([]model.WalkInData, error) {
	data := make([]model.WalkInData, 0)
	query := `SELECT COUNT(1) as total, cast(convert_tz(created_at, 'UTC', 'Asia/Manila') as date)  as date FROM client_log where created_at BETWEEN DATE_SUB(NOW(), INTERVAL 7 DAY) AND NOW() GROUP BY DAY(created_at), cast(convert_tz(created_at, 'UTC', 'Asia/Manila') as date) `
	err := repo.db.Select(&data, query)
	return data, err
}
func (repo *Dashboard)GetMonthlyWalkIns()([]model.WalkInData, error) {
	data := make([]model.WalkInData, 0)
	query := `SELECT COUNT(1) as total, cast(convert_tz(created_at, 'UTC', 'Asia/Manila') as date)   as date FROM client_log where created_at BETWEEN DATE_SUB(NOW(), INTERVAL 1 MONTH) AND NOW() GROUP BY DAY(created_at), cast(convert_tz(created_at, 'UTC', 'Asia/Manila') as date)`
	err := repo.db.Select(&data, query)
	return data, err
}
func (repo Dashboard)GetWeeklyCoachClients (coachId int) ([]model.CoachClient, error ) {
  
   data := make([]model.CoachClient, 0)
   query := `SELECT COUNT(1) as total, cast(convert_tz(created_at, 'UTC', 'Asia/Manila') as date)  
      as date FROM hired_coach where hired_coach.coach_id = ? and status_id = 3 and created_at BETWEEN DATE_SUB(NOW(), INTERVAL 7 DAY) AND NOW() 
      GROUP BY DAY(created_at), cast(convert_tz(created_at, 'UTC', 'Asia/Manila') as date) 
   `
   err := repo.db.Select(&data, query, coachId)
   return data, err
}

func (repo Dashboard)GetMonthlyCoachClients (coachId int) ([]model.CoachClient, error ) {
  data := make([]model.CoachClient, 0)
  query := `SELECT COUNT(1) as total, cast(convert_tz(created_at, 'UTC', 'Asia/Manila') as date)  
     as date FROM hired_coach where hired_coach.coach_id = ? and  status_id = 3 and created_at BETWEEN DATE_SUB(NOW(), INTERVAL 1 MONTH) AND NOW() 
     GROUP BY DAY(created_at), cast(convert_tz(created_at, 'UTC', 'Asia/Manila') as date) 
  `
  err := repo.db.Select(&data, query, coachId)
  return data, err
}
func (repo * Dashboard)GetClientDashboardData(clientId int) (model.ClientDashboardData,error) {
  data := model.ClientDashboardData{}
  query := `
  SELECT 
  (SELECT count(1) from package_request where client_id = ?) as packages,
  (SELECT count(1) from reservation where client_id = ?) as reservations,
  (SELECT count(1) from hired_coach where client_id = ?) as coach_appointments,
  (SELECT count(1) from membership_request where client_id = ?) as membership_requests,
  JSON_OBJECT(
        'walkIn', COALESCE((SELECT SUM(amount_paid) 
        FROM client_log  
        WHERE deleted_at is null and created_at BETWEEN DATE_SUB(NOW(), INTERVAL 1 YEAR) AND NOW() and client_id = ?), 0), 
    'membership',  COALESCE((SELECT  SUM(price) from subscription
        INNER JOIN membership_plan_snapshot on subscription.membership_plan_snapshot_id = membership_plan_snapshot.id
      where subscription.valid_until >= NOW() 
      and subscription.cancelled_at is NULL and client_id = ? and subscription.created_at BETWEEN DATE_SUB(NOW(), INTERVAL 1 YEAR) 
        AND NOW()),0),
        'package', COALESCE((SELECT SUM(price) from package_request 
        INNER JOIN package_snapshot on package_request.package_snapshot_id = package_snapshot_id 
        where client_id = ? and package_request.created_at BETWEEN DATE_SUB(NOW(), INTERVAL 1 YEAR) 
        AND NOW()),0)
      ) as annual_expenditures_breakdown,
        
        JSON_OBJECT(
        'walkIn', COALESCE((SELECT SUM(amount_paid) 
        FROM client_log  
        WHERE deleted_at is null and client_id = ? and created_at BETWEEN DATE_SUB(NOW(), INTERVAL 1 MONTH) AND NOW()), 0), 
    'membership',  COALESCE((SELECT  SUM(price) from subscription
        INNER JOIN membership_plan_snapshot on subscription.membership_plan_snapshot_id = membership_plan_snapshot.id
      where client_id = ? and subscription.valid_until >= NOW() 
      and subscription.cancelled_at is NULL and subscription.created_at BETWEEN DATE_SUB(NOW(), INTERVAL 1 MONTH) 
        AND NOW()),0),
        'package', COALESCE((SELECT SUM(price) from package_request 
        INNER JOIN package_snapshot on package_request.package_snapshot_id = package_snapshot_id 
        where client_id = ? and package_request.created_at BETWEEN DATE_SUB(NOW(), INTERVAL 1 MONTH) 
        AND NOW()),0)
      ) as monthly_expenditures_breakdown,
        JSON_OBJECT(	
        'walkIn', COALESCE((SELECT SUM(amount_paid) 
        FROM client_log  
        WHERE deleted_at is null and client_id = ? and  created_at BETWEEN DATE_SUB(NOW(), INTERVAL 7 DAY) AND NOW()), 0), 
    'membership',  COALESCE((SELECT  SUM(price) from subscription
        INNER JOIN membership_plan_snapshot on subscription.membership_plan_snapshot_id = membership_plan_snapshot.id
      where client_id = ? and  subscription.valid_until >= NOW() 
      and subscription.cancelled_at is NULL and subscription.created_at BETWEEN DATE_SUB(NOW(), INTERVAL 7 DAY) 
        AND NOW()),0),
        'package', COALESCE((SELECT SUM(price) from package_request 
        INNER JOIN package_snapshot on package_request.package_snapshot_id = package_snapshot_id 
        where client_id = ? and package_request.created_at BETWEEN DATE_SUB(NOW(), INTERVAL 7 DAY) 
        AND NOW()),0)
  ) as weekly_expenditures_breakdown
  `
  err := repo.db.Get(&data, query, clientId, clientId, clientId, clientId,clientId,clientId,clientId,clientId,clientId, clientId, clientId, clientId, clientId )
	return data, err

}

func (repo * Dashboard)GetClientWalkIns(clientId int) ([]model.ClientLog, error){
	logs := make([]model.ClientLog, 0)
	query := `SELECT client_log.id, client_log.client_id, client_log.amount_paid, client_log.is_member, JSON_OBJECT('publicId',client.public_id ,'id', client.id, 'givenName', client.given_name, 'middleName', client.middle_name, 'surname', client.surname, 'email', account.email)  as client,  convert_tz(client_log.created_at, 'UTC', 'Asia/Manila') as created_at from client_log
		INNER JOIN client on client_log.client_id = client.id
		INNER JOIN account on client.account_id = account.id
		where client_log.deleted_at is null and client_id = ?
		ORDER BY client_log.created_at DESC`

	err := repo.db.Select(&logs, query, clientId)
	return logs, err
}

func(repo Dashboard)GetCoachDashboardData (coachId int) (model.CoachDashboardData, error) {
  data := model.CoachDashboardData{}
  query := `
     SELECT 
    (SELECT COUNT(DISTINCT client_id) from hired_coach where coach_id = ?) as clients,
    (SELECT COUNT(1) from hired_coach where coach_id = ?) as appointments,
    (SELECT 
    SUM(coaching_rate_snapshot.price) 
    from hired_coach 
    inner join coaching_rate_snapshot 
    on hired_coach.rate_snapshot_id = coaching_rate_snapshot.id 
    where hired_coach.coach_id = ? and status_id = 3 ) as earnings
  `
  err := repo.db.Get(&data, query, coachId, coachId, coachId)
  return data, err
}

