package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/status"

	"github.com/jmoiron/sqlx"
)

type MembershipRequestRepository struct {
	db *sqlx.DB
}
func(repo *MembershipRequestRepository)NewRequest(request model.MembershipRequest) error {
	transaction, err := repo.db.Beginx()
	if err != nil {
		transaction.Rollback()
		return err
	}
	plan := model.MembershipPlan{}
	getErr := transaction.Get(&plan, "SELECT price, description, months from membership_plan where id = ?", request.MembershipPlanId)
	if getErr != nil {
		transaction.Rollback()
		return getErr
	}
	result , err := transaction.Exec("INSERT INTO membership_plan_snapshot(description,months, price) VALUES(?, ?, ?)", plan.Description, plan.Months, plan.Price)
	if err != nil {
		transaction.Rollback()
		return err
	}
	snapshotId, err := result.LastInsertId()
	_, err = transaction.Exec("INSERT INTO membership_request(client_id, membership_plan_id, status_id, membership_plan_snapshot_id, remarks )VALUES(?, ?, ?, ?, ?)", request.ClientId, request.MembershipPlanId, request.StatusId, snapshotId, "")
	if err != nil {
		transaction.Rollback()
		return err
	}
	transaction.Commit()
	return nil
}

func (repo MembershipRequestRepository)GetMembershipRequestsByClientId(clientId int) ([]model.MembershipRequest, error) {
	requests := make([]model.MembershipRequest, 0)
	err := repo.db.Select(&requests,`
	SELECT 
	mbr.id, mbr.client_id, mbr.membership_plan_id, mbr.status_id, mbrs.description as status, mbr.remarks,
	JSON_OBJECT('publicId', client.public_id, 'id', client.id, 'givenName', client.given_name, 'middleName', client.middle_name, 'surname', client.surname)  as client,
	JSON_OBJECT('id', mp.id, 'description', mp.description, 'months', mp.months, 'price', mp.price) as membership_plan,
    JSON_OBJECT('id', membership_plan_snapshot.id, 'description', membership_plan_snapshot.description, 'months', membership_plan_snapshot.months, 'price', membership_plan_snapshot.price) as membership_plan_snapshot,
    mbr.created_at
	FROM membership_request as mbr 
	INNER JOIN membership_request_status as mbrs on mbr.status_id = mbrs.id
    INNER JOIN membership_plan_snapshot on mbr.membership_plan_snapshot_id = membership_plan_snapshot.id
	INNER JOIN client on mbr.client_id = client.id
	INNER JOIN membership_plan as mp on mbr.membership_plan_id = mp.id
	WHERE mbr.client_id = ? ORDER BY mbr.updated_at
	`, clientId)
	return requests, err
}

func (repo MembershipRequestRepository)GetMembershipRequestById(id int) (model.MembershipRequest, error) {
	request := model.MembershipRequest{}
	err := repo.db.Get(&request,`
	SELECT 
	mbr.id, mbr.client_id, mbr.membership_plan_id, mbr.status_id, mbrs.description as status, mbr.remarks,
	JSON_OBJECT('publicId', client.public_id,'id', client.id, 'givenName', client.given_name, 'middleName', client.middle_name, 'surname', client.surname)  as client,
	JSON_OBJECT('id', mp.id, 'description', mp.description, 'months', mp.months, 'price', mp.price) as membership_plan,
    JSON_OBJECT('id', membership_plan_snapshot.id, 'description', membership_plan_snapshot.description, 'months', membership_plan_snapshot.months, 'price', membership_plan_snapshot.price) as membership_plan_snapshot,
    mbr.created_at
	FROM membership_request as mbr 
	INNER JOIN membership_request_status as mbrs on mbr.status_id = mbrs.id
    INNER JOIN membership_plan_snapshot on mbr.membership_plan_snapshot_id = membership_plan_snapshot.id
	INNER JOIN client on mbr.client_id = client.id
	INNER JOIN membership_plan as mp on mbr.membership_plan_id = mp.id
	WHERE mbr.id = ? ORDER BY mbr.updated_at
	`, id)
	return request, err
}
func (repo MembershipRequestRepository)GetMembershipRequests() ([]model.MembershipRequest, error) {
	requests := make([]model.MembershipRequest, 0)
	err := repo.db.Select(&requests,`
	SELECT 
	mbr.id, mbr.client_id, mbr.membership_plan_id, mbr.status_id, mbrs.description as status, mbr.remarks,
	JSON_OBJECT('publicId', client.public_id,'id', client.id, 'givenName', client.given_name, 'middleName', client.middle_name, 'surname', client.surname)  as client,
	JSON_OBJECT('id', mp.id, 'description', mp.description, 'months', mp.months, 'price', mp.price) as membership_plan,
    JSON_OBJECT('id', membership_plan_snapshot.id, 'description', membership_plan_snapshot.description, 'months', membership_plan_snapshot.months, 'price', membership_plan_snapshot.price) as membership_plan_snapshot,
    mbr.created_at
	FROM membership_request as mbr 
	INNER JOIN membership_request_status as mbrs on mbr.status_id = mbrs.id
    INNER JOIN membership_plan_snapshot on mbr.membership_plan_snapshot_id = membership_plan_snapshot.id
	INNER JOIN client on mbr.client_id = client.id
	INNER JOIN membership_plan as mp on mbr.membership_plan_id = mp.id
	ORDER BY mbr.updated_at;
	`)
	return requests, err
}

func(repo * MembershipRequestRepository) CancelMembershipRequest( id int, remarks string) error {
	_, err := repo.db.Exec(
		"UPDATE membership_request SET status_id = ?, remarks = ? where id = ? AND status_id != ?", 
		status.MembershipRequestStatusCancelled, remarks,  id, status.MembershipRequestStatusReceived)
	return err
}
func(repo * MembershipRequestRepository) ApproveMembershipRequest( id int, remarks string) error {
	_, err := repo.db.Exec(
		"UPDATE membership_request SET status_id = ?, remarks = ? where id = ? AND status_id = ?", 
		status.MembershipRequestStatusApproved, remarks,  id, status.MembershipRequestStatusPending)
	return err
}

func (repo * MembershipRequestRepository)MarkAsReceived(id int, remarks string) error {
	_, err := repo.db.Exec(
		"UPDATE membership_request SET status_id = ?, remarks = ? where id = ? AND status_id = ?", 
		status.MembershipRequestStatusReceived, remarks,  id, status.MembershipRequestStatusApproved)
	return err
}
func NewMembershipRequestRepository() MembershipRequestRepository{
	return MembershipRequestRepository {
		db: db.GetConnection(),
	}
}


