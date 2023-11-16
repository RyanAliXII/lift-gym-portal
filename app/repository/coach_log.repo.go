package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)

type CoachLogRepository struct {
	db *sqlx.DB
}
func NewCoachLogRepository() CoachLogRepository {
	return CoachLogRepository{
		db: db.GetConnection(),
	}
}
func (repo * CoachLogRepository)NewLog(log model.CoachLog) error{
	_, err := repo.db.Exec("INSERT INTO coach_log(coach_id) VALUES(?)", log.CoachId)
	return err
}

func (repo * CoachLogRepository)UpdateLog(log model.CoachLog) error{
	_, err := repo.db.Exec("UPDATE coach_log set coach_id = ? WHERE id = ?" ,log.CoachId,  log.Id)
	return err
}
func (repo * CoachLogRepository)GetLogs() ([]model.CoachLog, error){
	logs := make([]model.CoachLog, 0)
	query := `SELECT coach_log.id, 
	coach_log.coach_id, 
	(case when coach_log.logged_out_at is null then false else true end) as is_logged_out,
	(case when coach_log.logged_out_at is null then '' else  convert_tz(coach_log.logged_out_at, 'UTC', 'Asia/Manila')  end) as logged_out_at,
	JSON_OBJECT('publicId',coach.public_id ,'id', coach.id, 'givenName', coach.given_name, 'middleName', coach.middle_name, 'surname', coach.surname, 'email', account.email)  as coach,  convert_tz(coach_log.created_at, 'UTC', 'Asia/Manila') as created_at from coach_log
		INNER JOIN coach on  coach_log.coach_id = coach.id
		INNER JOIN account on coach.account_id = account.id
		where coach_log.deleted_at is null
		ORDER BY coach_log.created_at DESC`

	err := repo.db.Select(&logs, query)
	return logs, err
}

func (repo * CoachLogRepository)DeleteLog(id int) error{
	_, err := repo.db.Exec("UPDATE coach_log set deleted_at = now() where id = ?", id)
	return err
}

func (repo * CoachLogRepository)LogoutCoach(id int) error{
	_, err := repo.db.Exec("UPDATE coach_log set logged_out_at  = now() where id = ?", id)
	return err
}
