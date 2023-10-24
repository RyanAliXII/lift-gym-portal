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

func (repo *RoleRepository) GetRoles() ([]model.Role, error) {

	roles := make([]model.Role,0)
	query := `
	SELECT role.id, name, CONCAT('[',GROUP_CONCAT('"',permission.value,'"'),']') as permissions from role
	INNER JOIN permission on role.id = permission.role_id GROUP BY role.id ORDER BY role.updated_at DESC
	`
	err := repo.db.Select(&roles, query)
	return roles, err
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


func (repo *RoleRepository) UpdateRole(role model.Role) error {
	transaction, err := repo.db.Begin()
	if err != nil {
		transaction.Rollback()
		return err
	}
	_, err = transaction.Exec("UPDATE role set name = ? where id = ?", role.Name, role.Id)
	if err != nil {
		transaction.Rollback()
		return err
	}

	_, err = transaction.Exec("DELETE FROM permission where role_id = ?", role.Id)
	if err != nil {
		transaction.Rollback()
		return err
	}
	records := make([]goqu.Record, 0)

	for _, permission := range role.Permissions {
		records = append(records, goqu.Record{
			"role_id": role.Id,
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
func (repo * RoleRepository)GetRoleByUserId(userId int)(model.Role, error){
	role := model.Role{}
	err := repo.db.Get(&role, `
	SELECT role.id, role.name, CONCAT('[',GROUP_CONCAT('"',permission.value,'"'),']') as permissions from user 
	INNER JOIN role on  user.role_id = role.id
	INNER JOIN permission on role.id =permission.role_id 
	where user.id = ?
	GROUP BY user.id, role.id;
	`, userId)
	return role, err
}