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
    (SELECT COUNT(1) FROM client) AS clients, 
    (SELECT COUNT(1) 
     FROM client 
     INNER JOIN subscription ON subscription.client_id = client.id  
     WHERE subscription.valid_until >= NOW() AND subscription.cancelled_at IS NULL) AS members,
    (
     COALESCE((SELECT SUM(amount_paid) 
     FROM client_log  
     WHERE created_at BETWEEN DATE_SUB(NOW(), INTERVAL 30 DAY) AND NOW()), 0) 
     + 
     COALESCE((SELECT  SUM(price) from subscription
     INNER JOIN membership_plan_snapshot on subscription.membership_plan_snapshot_id = membership_plan_snapshot.id
	  where subscription.valid_until >= NOW() 
	  and subscription.cancelled_at is NULL and subscription.created_at BETWEEN DATE_SUB(NOW(), INTERVAL 30 DAY) 
      AND NOW()),0)
      + 
      COALESCE((
      SELECT SUM(price) from package_request 
      INNER JOIN package_snapshot on package_request.package_snapshot_id = package_snapshot_id 
      where package_request.created_at BETWEEN DATE_SUB(NOW(), INTERVAL 30 DAY) 
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
      WHERE created_at BETWEEN DATE_SUB(NOW(), INTERVAL 30 DAY) AND NOW()), 0), 
	 'membership',  COALESCE((SELECT  SUM(price) from subscription
      INNER JOIN membership_plan_snapshot on subscription.membership_plan_snapshot_id = membership_plan_snapshot.id
	  where subscription.valid_until >= NOW() 
	  and subscription.cancelled_at is NULL and subscription.created_at BETWEEN DATE_SUB(NOW(), INTERVAL 30 DAY) 
      AND NOW()),0),
      'package', COALESCE((SELECT SUM(price) from package_request 
      INNER JOIN package_snapshot on package_request.package_snapshot_id = package_snapshot_id 
      where package_request.created_at BETWEEN DATE_SUB(NOW(), INTERVAL 30 DAY) 
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