package repository

import (
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
		insertClientQuery := `INSERT INTO client(given_name, middle_name, surname, date_of_birth, address, mobile_number, emergency_contact, account_id) VALUES(?, ? ,? , ?, ? , ?, ?, ?)`
		_, insertClientErr := transaction.Exec(insertClientQuery, client.GivenName, client.MiddleName, client.Surname, client.DateOfBirth, client.Address, client.MobileNumber, client.EmergencyContact, accountId )
		if insertClientErr != nil {
			transaction.Rollback()
			return insertClientErr
		}
		transaction.Commit()
		return nil
}
func (repo * ClientRepository)Get()([]model.Client, error) {
	clients := make([]model.Client , 0)
	selectQuery := `SELECT client.id, client.given_name, client.middle_name, client.surname, client.date_of_birth, client.address, client.emergency_contact,client.mobile_number, account.email, account.id as account_id from client
	INNER JOIN account on client.account_id = account.id ORDER BY client.updated_at DESC`
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
func (repo *ClientRepository) GetClientByEmail(email string) (model.Client, error) {
	client := model.Client{}
	getQuery := `SELECT client.id, client.given_name, client.middle_name, client.surname, client.date_of_birth, client.emergency_contact,client.mobile_number, account.email from client
		INNER JOIN account on client.account_id = account.id where UPPER(account.email) = UPPER(?) LIMIT 1`
	getErr := repo.db.Get(&client, getQuery, email)
	return client, getErr 
}

func (repo * ClientRepository) GetClientById (id int) (model.Client, error) {
	client := model.Client{}
	getQuery := `SELECT client.id, client.given_name, client.middle_name, client.surname, client.date_of_birth, client.address, client.emergency_contact,client.mobile_number, account.email, account.id as account_id from client
	INNER JOIN account on client.account_id = account.id where client.id = ? LIMIT 1`
	getErr := repo.db.Get(&client, getQuery, id)
	return client, getErr
}
func NewClientRepository () ClientRepository{
	return ClientRepository{
		db: db.GetConnection(),
	}
}
