package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)

type PackageRequestRepository struct {
	db *sqlx.DB
}

func (repo *PackageRequestRepository) NewPackageRequest(pkgRequest model.PackageRequest)(error){
	query := `INSERT INTO package_request(client_id, package_id, status_id) VALUES (?, ?, ?)`
	_, err:= repo.db.Exec(query, pkgRequest.ClientId, pkgRequest.PackageId, pkgRequest.StatusId)
	return err
}


func NewPackageRequestRepository() PackageRequestRepository{
	return PackageRequestRepository{
		db: db.GetConnection(),
	}
}