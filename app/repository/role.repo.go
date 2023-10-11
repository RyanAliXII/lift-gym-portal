package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/jmoiron/sqlx"
)

type RoleRepository struct {
	db * sqlx.DB
}
func NewRoleRepository() RoleRepository{
	return RoleRepository{
		db: db.GetConnection(),
	}
}
func (repo *RoleRepository) NewRole(role model.Role) error {
	transaction, err := repo.db.Begin()
	if err != nil {
		transaction.Rollback()
		return err
	}
	result, err := transaction.Exec("INSERT into role (name) VALUES(?)", role.Name)
	if err != nil {
		transaction.Rollback()
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		transaction.Rollback()
		return err
	}
	records := make([]goqu.Record, 0)

	for _, permission := range role.Permissions {
		records = append(records, goqu.Record{
			"role_id": id,
			"value": permission,
		})
	}
	dialect := goqu.Dialect("mysql")
	ds := dialect.Insert(goqu.T("permission")).Rows(records).Prepared(true)
	query, args, _ := ds.ToSQL()
	_, err = transaction.Exec(query, args...)
	if err != nil{
	     transaction.Rollback()
		 return err
	}
	transaction.Commit()
	return nil;
}
