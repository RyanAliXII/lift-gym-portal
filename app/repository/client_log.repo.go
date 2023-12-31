package repository

import (
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)

type ClientLogRepository struct {
	db *sqlx.DB
}
func NewClientLogRepository() ClientLogRepository {
	return ClientLogRepository{
		db: db.GetConnection(),
	}
}
func (repo * ClientLogRepository)NewLog(log model.ClientLog) error{
	_, err := repo.db.Exec("INSERT INTO client_log(client_id, amount_paid, is_member) VALUES(?, ?, ?)", log.ClientId, log.AmountPaid, log.IsMember)
	return err
}

func (repo * ClientLogRepository)UpdateLog(log model.ClientLog) error{
	_, err := repo.db.Exec("UPDATE client_log set client_id = ?, amount_paid = ?, is_member = ? WHERE id = ?" ,log.ClientId, log.AmountPaid, log.IsMember, log.Id)
	return err
}
func (repo * ClientLogRepository)GetLogs()([]model.ClientLog, error){
	logs := make([]model.ClientLog, 0)
	query := `SELECT client_log.id, 
	client_log.client_id, 
	client_log.amount_paid, 
	client_log.is_member, 
	(case when client_log.logged_out_at is null then false else true end) as is_logged_out,
	(case when client_log.logged_out_at is null then '' else  convert_tz(client_log.logged_out_at, 'UTC', 'Asia/Manila')  end) as logged_out_at,
	JSON_OBJECT('publicId',client.public_id ,'id', client.id, 'givenName', client.given_name, 'middleName', client.middle_name, 'surname', client.surname, 'email', account.email)  as client,  convert_tz(client_log.created_at, 'UTC', 'Asia/Manila') as created_at from client_log
		INNER JOIN client on client_log.client_id = client.id
		INNER JOIN account on client.account_id = account.id
		where client_log.deleted_at is null
		ORDER BY client_log.created_at DESC`

	err := repo.db.Select(&logs, query)
	return logs, err
}

func (repo * ClientLogRepository)DeleteLog(id int) error{
	_, err := repo.db.Exec("UPDATE client_log set deleted_at = now() where id = ?", id)
	return err
}

func (repo * ClientLogRepository)LogoutClient(id int) error{
	_, err := repo.db.Exec("UPDATE client_log set logged_out_at  = now() where id = ?", id)
	return err
}
