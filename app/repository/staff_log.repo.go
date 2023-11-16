package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)

type StaffLogRepository struct {
	db *sqlx.DB
}
func NewStaffLogRepository() StaffLogRepository {
	return StaffLogRepository{
		db: db.GetConnection(),
	}
}
func (repo * StaffLogRepository)NewLog(log model.StaffLog) error{
	_, err := repo.db.Exec("INSERT INTO staff_log(staff_id) VALUES(?)", log.StaffId)
	return err
}

func (repo * StaffLogRepository)UpdateLog(log model.StaffLog) error{
	_, err := repo.db.Exec("UPDATE staff_log set staff_id = ? WHERE id = ?" ,log.StaffId,  log.Id)
	return err
}
func (repo * StaffLogRepository)GetLogs() ([]model.StaffLog, error){
	logs := make([]model.StaffLog, 0)
	query := `SELECT staff_log.id, 
	staff_log.staff_id, 
	(case when staff_log.logged_out_at is null then false else true end) as is_logged_out,
	(case when staff_log.logged_out_at is null then '' else  convert_tz(staff_log.logged_out_at, 'UTC', 'Asia/Manila')  end) as logged_out_at,
	JSON_OBJECT('publicId', staff.public_id ,'id', staff.id, 'givenName', staff.given_name, 'middleName', staff.middle_name, 'surname', staff.surname, 'email', account.email)  as staff,  convert_tz(staff_log.created_at, 'UTC', 'Asia/Manila') as created_at from staff_log
		INNER JOIN user as staff on staff_log.staff_id = staff.id
		INNER JOIN account on staff.account_id = account.id
		where staff_log.deleted_at is null
		ORDER BY staff_log.created_at DESC`

	err := repo.db.Select(&logs, query)
	return logs, err
}

func (repo * StaffLogRepository)DeleteLog(id int) error{
	_, err := repo.db.Exec("UPDATE staff_log set deleted_at = now() where id = ?", id)
	return err
}

func (repo * StaffLogRepository)LogoutStaff(id int) error{
	_, err := repo.db.Exec("UPDATE staff_log set logged_out_at  = now() where id = ?", id)
	return err
}
