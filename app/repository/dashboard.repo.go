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
     WHERE created_at BETWEEN DATE_SUB(NOW(), INTERVAL 1 YEAR) AND NOW()), 0) 
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
     WHERE created_at BETWEEN DATE_SUB(NOW(), INTERVAL 7 DAY) AND NOW()), 0) 
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
      WHERE created_at BETWEEN DATE_SUB(NOW(), INTERVAL 1 YEAR) AND NOW()), 0), 
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
      WHERE created_at BETWEEN DATE_SUB(NOW(), INTERVAL 1 MONTH) AND NOW()), 0), 
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

