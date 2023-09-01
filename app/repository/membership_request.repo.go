package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)

type MembershipRequestRepository struct {
	db *sqlx.DB
}
func(repo *MembershipRequestRepository)NewRequest(request model.MembershipRequest) error {
	_, err := repo.db.Exec("INSERT INTO membership_request(client_id, membership_plan_id, status_id)VALUES(?, ?, ?)", request.ClientId, request.MembershipPlanId, request.StatusId)
	return err
}
func NewMembershipRequestRepository() MembershipRequestRepository{
	return MembershipRequestRepository {
		db: db.GetConnection(),
	}
}


