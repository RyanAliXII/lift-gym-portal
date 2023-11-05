package repository

import (
	"fmt"
	"lift-fitness-gym/app/db"
	"lift-fitness-gym/app/model"

	"github.com/jmoiron/sqlx"
)

type ClientRepository struct {
	db *sqlx.DB
}
func (repo *ClientRepository) New(client model.Client)(error) {
		transaction, transactErr := repo.db.Begin()
		if transactErr != nil {
			transaction.Rollback()
			return transactErr
		}
		result, insertAccountErr := transaction.Exec("INSERT INTO account (email, password)  VALUES (?, ?)", client.Email, client.Password)
		if insertAccountErr != nil {
			transaction.Rollback()
			return insertAccountErr
		}
		accountId, lastInsertedIdErr := result.LastInsertId()
		if lastInsertedIdErr != nil {
			transaction.Rollback()
			return lastInsertedIdErr
		}
		insertClientQuery := `INSERT INTO client(given_name, middle_name, surname, date_of_birth, address, mobile_number, emergency_contact, account_id, gender) VALUES(?, ? ,? , ?, ? , ?, ?, ?, ?)`
		_, insertClientErr := transaction.Exec(insertClientQuery, client.GivenName, client.MiddleName, client.Surname, client.DateOfBirth, client.Address, client.MobileNumber, client.EmergencyContact, accountId, client.Gender)
		if insertClientErr != nil {
			transaction.Rollback()
			return insertClientErr
		}
		transaction.Commit()
		return nil
}
func (repo * ClientRepository)Get()([]model.Client, error) {
	clients := make([]model.Client , 0)
	selectQuery := `SELECT client.id, client.given_name, client.middle_name,client.surname,client.public_id, client.date_of_birth, client.gender, client.address, client.emergency_contact,client.mobile_number, account.email, account.id as account_id, (case when verified_at is null then false else true end) as is_verified, (case when subscription.id then true else false end) as is_member from client
	INNER JOIN account on client.account_id = account.id 
	LEFT JOIN subscription on subscription.client_id = client.id
	AND subscription.valid_until >= NOW() and subscription.cancelled_at is NULL
	where client.deleted_at is NULL
	ORDER BY client.updated_at DESC`
	selectErr := repo.db.Select(&clients, selectQuery)
	return clients, selectErr 
}

func (repo * ClientRepository)Search(keyword string)([]model.Client, error) {
	clients := make([]model.Client , 0)
	keywordLike := fmt.Sprintf("%s%s%s","%", keyword, "%")
	selectQuery := `SELECT client.id, client.given_name, client.middle_name, client.surname,client.public_id, client.date_of_birth, client.gender, client.address, client.emergency_contact,client.mobile_number, account.email, account.id as account_id, (case when verified_at is null then false else true end) as is_verified, (case when subscription.id then true else false end) as is_member from client
	INNER JOIN account on client.account_id = account.id 
	LEFT JOIN subscription on subscription.client_id = client.id
	AND subscription.valid_until >= NOW() and subscription.cancelled_at is NULL
    where (client.given_name LIKE ? OR client.middle_name LIKE ? OR client.surname LIKE ? OR client.mobile_number LIKE ? OR account.email LIKE ? OR client.public_id LIKE ?) and client.deleted_at is null
	ORDER BY client.updated_at DESC LIMIT 50`
	selectErr := repo.db.Select(&clients, selectQuery, keywordLike, keywordLike, keywordLike, keywordLike, keywordLike, keywordLike)
	return clients, selectErr 
}
func (repo * ClientRepository)GetById(id int)(model.Client, error) {
	clients := model.Client{}
	selectQuery := `SELECT client.id, client.given_name, client.middle_name, client.surname, client.public_id, client.date_of_birth, client.gender, client.address, client.emergency_contact,client.mobile_number, account.email, account.id as account_id, (case when verified_at is null then false else true end) as is_verified, (case when subscription.id then true else false end) as is_member from client
	INNER JOIN account on client.account_id = account.id 
	LEFT JOIN subscription on subscription.client_id = client.id
	AND subscription.valid_until >= NOW() and subscription.cancelled_at is NULL
	where client.id = ? and client.deleted_at is null
	ORDER BY client.updated_at DESC;`
	getErr := repo.db.Get(&clients, selectQuery , id)
	return clients, getErr 
}
func (repo * ClientRepository)GetByIdWithPassword(id int)(model.Client, error) {
	clients := model.Client{}
	selectQuery := `SELECT client.id, client.given_name, client.middle_name, client.surname, client.public_id, client.date_of_birth, client.gender, client.address, client.emergency_contact,client.mobile_number, account.email, account.password, account.id as account_id, (case when verified_at is null then false else true end ) as is_verified from client
	INNER JOIN account on client.account_id = account.id where client.id = ? and deleted_at is null ORDER BY client.updated_at DESC LIMIT 1;`
	getErr := repo.db.Get(&clients, selectQuery , id)
	return clients, getErr 
}
func (repo * ClientRepository)GetUnsubscribed()([]model.Client, error) {
	clients := make([]model.Client , 0)
	selectQuery := `SELECT client.id, client.given_name, client.middle_name, client.surname, client.public_id, client.date_of_birth, client.gender, client.address, client.emergency_contact, client.mobile_number, account.email, account.id as account_id
	FROM client
	INNER JOIN account ON client.account_id = account.id
	WHERE client.id NOT IN (
		SELECT subscription.client_id
		FROM subscription
		WHERE subscription.valid_until >= NOW() AND subscription.cancelled_at IS NULL
	) and client.deleted_at is null
	ORDER BY client.updated_at DESC;`
	selectErr := repo.db.Select(&clients, selectQuery)
	return clients, selectErr 
}

func (repo * ClientRepository)UpdatePassword (newPassword string, clientId int )(error){
	client, getClientErr := repo.GetClientById(clientId)
	if getClientErr != nil {
		return getClientErr
	}
	updateQuery := "UPDATE account SET password = ? WHERE id = ?"
	_, updateErr := repo.db.Exec(updateQuery, newPassword, client.AccountId)
	if updateErr != nil {
		return updateErr
	}
	return nil
} 

func (repo  ClientRepository) Update(client model.Client) (error) {
	dbClient, getClientErr := repo.GetClientById(client.Id)
	if getClientErr != nil {
		return getClientErr
	}
	transaction, transactErr := repo.db.Begin()
	if transactErr != nil {
		transaction.Rollback()
		return transactErr
	}
	updateClientQuery := `UPDATE client SET given_name = ?, middle_name = ?, surname = ?, date_of_birth = ?, address = ?, mobile_number = ?, emergency_contact = ?, gender = ? where id = ?`
	_, updateClientErr := transaction.Exec(updateClientQuery, client.GivenName,client.MiddleName, client.Surname, client.DateOfBirth, client.Address, client.MobileNumber, client.EmergencyContact, client.Gender, client.Id)
	if updateClientErr != nil {
		transaction.Rollback()
		return updateClientErr
	}
	updateAccountQuery := `UPDATE account SET email = ? WHERE id = ?`
	_, updateAccountErr := transaction.Exec(updateAccountQuery, client.Email, dbClient.AccountId)
	if updateAccountErr != nil {
		transaction.Rollback()
		return updateAccountErr
	}
	transaction.Commit()
	return nil
}
func (repo *ClientRepository) GetClientByEmail(email string) (model.Client, error) {
	client := model.Client{}
	getQuery := `SELECT client.id, client.given_name, client.middle_name, client.surname, client.date_of_birth, client.public_id, client.emergency_contact,client.mobile_number, account.email from client
		INNER JOIN account on client.account_id = account.id where UPPER(account.email) = UPPER(?) LIMIT 1`
	getErr := repo.db.Get(&client, getQuery, email)
	return client, getErr 
}
func (repo *ClientRepository) MarkAsVerified(id int ) error{
	_, err := repo.db.Exec("UPDATE client set verified_at = NOW() where id = ?", id)
	return err
}

func (repo * ClientRepository)Delete(id int)(error) {
	_, err := repo.db.Exec("UPDATE client set deleted_at = NOW() where id = ?", id)
	return err
}

func (repo * ClientRepository) GetClientById (id int) (model.Client, error) {
	client := model.Client{}
	getQuery := `SELECT client.id, client.given_name, client.middle_name, client.surname, client.public_id, client.date_of_birth, client.gender, client.address, client.emergency_contact,client.mobile_number, account.email, account.id as account_id from client
	INNER JOIN account on client.account_id = account.id where client.id = ? and client.deleted_at is null LIMIT 1`
	getErr := repo.db.Get(&client, getQuery, id)
	return client, getErr
}
func (repo * ClientRepository)UpdateMobileNumberOnce(clientId int, mobileNumber string) error {
	_, err := repo.db.Exec("UPDATE client set mobile_number = ? WHERE LENGTH(mobile_number) = 0 AND id = ?", mobileNumber, clientId)
	return err
}
func (repo * ClientRepository)UpdateEmergencyContactOnce(clientId int, mobileNumber string) error {
	_, err := repo.db.Exec("UPDATE client set emergency_contact = ? WHERE LENGTH(emergency_contact) = 0 AND id = ?", mobileNumber, clientId)
	return err
}
func (repo * ClientRepository)UpdateAddressOnce(clientId int, address string) error {
	_, err := repo.db.Exec("UPDATE client set address = ? WHERE LENGTH(address) = 0 AND id = ?", address, clientId)
	return err
}
func NewClientRepository () ClientRepository{
	return ClientRepository{
		db: db.GetConnection(),
	}
}
