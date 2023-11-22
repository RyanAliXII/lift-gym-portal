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
	SELECT role.id, name, (case when COALESCE(user_role.count, 0) > 0 then false else true end) as is_deletable, 
	COALESCE(JSON_ARRAYAGG(permission.value),'[]') as permissions
	FROM role
	INNER JOIN permission ON role.id = permission.role_id
	LEFT JOIN (
		SELECT role_id, COUNT(1) as count
		FROM user
		WHERE user.deleted_at IS NULL
		GROUP BY role_id
	) AS user_role ON role.id = user_role.role_id
	where deleted_at is null
	GROUP BY role.id, user_role.role_id, user_role.count
	ORDER BY role.updated_at DESC;
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

func (repo * RoleRepository)Delete(roleId int )error {
	role := model.Role{}
	//check if role is deletable before deleting.
	query := `
	SELECT role.id, name, (case when COALESCE(user_role.count, 0) > 0 then false else true end) as is_deletable, 
	COALESCE(JSON_ARRAYAGG(permission.value),'[]') as permissions
	FROM role
	INNER JOIN permission ON role.id = permission.role_id
	LEFT JOIN (
		SELECT role_id, COUNT(1) as count
		FROM user
		WHERE user.deleted_at IS NULL
		GROUP BY role_id
	) AS user_role ON role.id = user_role.role_id
	where (case when COALESCE(user_role.count, 0) > 0 then false else true end) = true and role.id = ?
	GROUP BY role.id, user_role.role_id, user_role.count
	ORDER BY role.updated_at DESC LIMIT 1`
	// this will throw error if row is empty
	err := repo.db.Get(&role, query, roleId)
	if err != nil {
		return err 
	}
	_, err = repo.db.Exec("Update role set deleted_at = now() where id = ?", role.Id)
	return err
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
func (repo * RoleRepository)GetRoleById(roleId int)(model.Role, error){
	role := model.Role{}
	err := repo.db.Get(&role, `
	SELECT role.id, name, (case when COALESCE(user_role.count, 0) > 0 then false else true end) as is_deletable, 
	COALESCE(JSON_ARRAYAGG(permission.value),'[]') as permissions
	FROM role
	INNER JOIN permission ON role.id = permission.role_id
	LEFT JOIN (
		SELECT role_id, COUNT(1) as count
		FROM user
		WHERE user.deleted_at IS NULL
		GROUP BY role_id
	) AS user_role ON role.id = user_role.role_id
	where deleted_at is null and role.id = ?
	GROUP BY role.id, user_role.role_id, user_role.count
	ORDER BY role.updated_at DESC LIMIT 1;
	`, roleId)
	return role, err
}