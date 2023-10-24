package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/status"

	"github.com/jmoiron/sqlx"
)

type PackageRepository struct {
	db * sqlx.DB
}


func (repo * PackageRepository) NewPackage (pkg model.Package) error{
	query := `INSERT INTO package(description, price)VALUES(?, ?)`
	_, insertErr := repo.db.Exec(query, pkg.Description, pkg.Price)

	return insertErr
}
func (repo * PackageRepository) UpdatePackage (pkg model.Package) error{
	query := `UPDATE package set description = ?, price = ? where id = ?`
	_, updateErr := repo.db.Exec(query, pkg.Description, pkg.Price, pkg.Id)
	return updateErr
}
func (repo * PackageRepository)GetPackages()([]model.Package, error){
	pkgs := make([]model.Package, 0)
	query := `SELECT id, description, price from package where deleted_at is null ORDER by updated_at DESC`
	selectErr := repo.db.Select(&pkgs, query)
	return pkgs, selectErr
}
func (repo * PackageRepository)DeletePackage(id int)(error){
	_, err := repo.db.Exec("UPDATE package set deleted_at = now() where id = ?", id)
	return err
}
func (repo * PackageRepository)GetUnrequestedPackageOfClient(clientId int)([]model.Package, error){
	pkgs := make([]model.Package, 0)
	err := repo.db.Select(&pkgs, `SELECT id, description, price FROM package 
	where package.id NOT IN(SELECT package_request.package_id FROM package_request where package_request.client_id = ?  AND (package_request.status_id = ? OR package_request.status_id = ?)) AND package.deleted_at is null;`, 
	clientId, status.PackageRequestStatusPending, status.PackageRequestStatusApproved)
	return pkgs, err
}

func NewPackageRepository () PackageRepository {
	return PackageRepository{
		db: db.GetConnection(),
	}
}