package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

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
	query := `SELECT id, description, price from package ORDER by updated_at DESC`
	selectErr := repo.db.Select(&pkgs, query)
	return pkgs, selectErr
}

func NewPackageRepository () PackageRepository {
	return PackageRepository{
		db: db.GetConnection(),
	}
}