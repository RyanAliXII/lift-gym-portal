package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/status"

	"github.com/jmoiron/sqlx"
)

type PackageRequestRepository struct {
	db *sqlx.DB
}

func (repo *PackageRequestRepository) NewPackageRequest(pkgRequest model.PackageRequest)(error){
	transaction, err := repo.db.Beginx()
	if err != nil {
		transaction.Rollback()
		return err
	}
	pkg := model.Package{}
	err = transaction.Get(&pkg, "SELECT description, price from package where id = ?", pkgRequest.PackageId)
	if err != nil {
		transaction.Rollback()
		return err
	}
	result, err := transaction.Exec("INSERT INTO package_snapshot (description, price) VALUES (?, ?)", pkg.Description, pkg.Price)
	if err != nil {
		transaction.Rollback()
		return err
	}
	snapshotId, err := result.LastInsertId()
	if err != nil {
		transaction.Rollback()
		return err
	}
	query := `INSERT INTO package_request(client_id, package_id, status_id, package_snapshot_id, remarks) VALUES (?, ?, ?, ?, "")`
	_, err = transaction.Exec(query, pkgRequest.ClientId, pkgRequest.PackageId, pkgRequest.StatusId, snapshotId)
	if err != nil {
		transaction.Rollback()
		return err
	}
	transaction.Commit()
	return nil
}
func (repo *PackageRequestRepository) NewPackageRequestWalkIn(pkgRequest model.PackageRequest)(error){
	transaction, err := repo.db.Beginx()
	if err != nil {
		transaction.Rollback()
		return err
	}
	pkg := model.Package{}
	err = transaction.Get(&pkg, "SELECT description, price from package where id = ?", pkgRequest.PackageId)
	if err != nil {
		transaction.Rollback()
		return err
	}
	result, err := transaction.Exec("INSERT INTO package_snapshot (description, price) VALUES (?, ?)", pkg.Description, pkg.Price)
	if err != nil {
		transaction.Rollback()
		return err
	}
	snapshotId, err := result.LastInsertId()
	if err != nil {
		transaction.Rollback()
		return err
	}
	query := `INSERT INTO package_request(client_id, package_id, status_id, package_snapshot_id, remarks) VALUES (?, ?, ?, ?, "")`
	_, err = transaction.Exec(query, pkgRequest.ClientId, pkgRequest.PackageId, pkgRequest.StatusId, snapshotId)
	if err != nil {
		transaction.Rollback()
		return err
	}
	transaction.Commit()
	return nil
}
func (repo *PackageRequestRepository) GetPackageRequestsByClientId (clientId int)([]model.PackageRequest, error) {	
	pkgRequests := make([]model.PackageRequest, 0)
	query := `SELECT 
	pkgr.id, pkgr.client_id, pkgr.package_id, pkgr.status_id, pkgrs.description as status, pkgr.remarks,
	JSON_OBJECT('id', client.id, 'givenName', client.given_name, 'middleName', client.middle_name, 'surname', client.surname)  as client,
	JSON_OBJECT('id', pkg.id, 'description', pkg.description, 'price', pkg.price) as package, 
	JSON_OBJECT('id', package_snapshot.id, 'description', package_snapshot.description, 'price', package_snapshot.price) as package_snapshot, pkgr.created_at
	FROM package_request as pkgr
	INNER JOIN package_request_status as pkgrs on pkgr.status_id = pkgrs.id
	INNER JOIN client on pkgr.client_id = client.id
	INNER JOIN package as pkg on pkgr.package_id = pkg.id
	INNER JOIN package_snapshot on pkgr.package_snapshot_id = package_snapshot.id
	where pkgr.client_id = ?
	ORDER BY pkgr.updated_at DESC`
	err := repo.db.Select(&pkgRequests, query, clientId)
	return pkgRequests, err
}
func (repo *PackageRequestRepository) GetPackageRequests()([]model.PackageRequest, error) {	
	pkgRequests := make([]model.PackageRequest, 0)
	query := `SELECT 
	pkgr.id, pkgr.client_id, pkgr.package_id, pkgr.status_id, pkgrs.description as status, pkgr.remarks,
	JSON_OBJECT('id', client.id, 'givenName', client.given_name, 'middleName', client.middle_name, 'surname', client.surname)  as client,
	JSON_OBJECT('id', pkg.id, 'description', pkg.description, 'price', pkg.price) as package, 
	JSON_OBJECT('id', package_snapshot.id, 'description', package_snapshot.description, 'price', package_snapshot.price) as package_snapshot, pkgr.created_at
	FROM package_request as pkgr
	INNER JOIN package_request_status as pkgrs on pkgr.status_id = pkgrs.id
	INNER JOIN client on pkgr.client_id = client.id
	INNER JOIN package as pkg on pkgr.package_id = pkg.id
	INNER JOIN package_snapshot on pkgr.package_snapshot_id = package_snapshot.id
	ORDER BY pkgr.updated_at DESC`
	err := repo.db.Select(&pkgRequests, query)
	return pkgRequests, err
}
func (repo * PackageRequestRepository)CancelPackageRequest(id int, remarks string) error {
	query := `UPDATE package_request set status_id = ?, remarks = ? where id = ? and status_id != ?`
	_, err := repo.db.Exec(query, status.PackageRequestStatusCancelled, remarks, id, status.PackageRequestStatusReceived)
	return err
}
func (repo * PackageRequestRepository)ApprovePackageRequest(id int, remarks string) error {
	query := `UPDATE package_request set status_id = ?, remarks = ? where id = ? and status_id = ?`
	_, err := repo.db.Exec(query, status.PackageRequestStatusApproved, remarks, id, status.PackageRequestStatusPending)
	return err
}
func (repo * PackageRequestRepository)MarkAsReceivedPackageRequest(id int, remarks string) error {
	query := `UPDATE package_request set status_id = ?, remarks = ? where id = ? and status_id = ?`
	_, err := repo.db.Exec(query, status.PackageRequestStatusReceived, remarks, id, status.PackageRequestStatusApproved)
	return err
}

func NewPackageRequestRepository() PackageRequestRepository{
	return PackageRequestRepository{
		db: db.GetConnection(),
	}
}