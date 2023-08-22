package db

import (
	"os"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)


func CreateRootAccount() {
	givenName := os.Getenv("ROOT_USER_GIVEN_NAME")
	middleName := os.Getenv("ROOT_USER_MIDDLE_NAME")
	surname := os.Getenv("ROOT_USER_SURNAME")
	email := os.Getenv("ROOT_USER_EMAIL")
	password := os.Getenv("ROOT_USER_PASSWORD")
	if len(givenName) == 0 || len(middleName) == 0 || len(surname) == 0 || len(email) == 0 || len(password) == 0{
		logger.Warn("One or more fields are empty.", zap.String("giveName", givenName), zap.String("middleName", middleName), 
		zap.String("surname", surname), zap.String("email", email), zap.String("password", password), )
		return
	}
	recordCount := 0
	db.Get(&recordCount,`SELECT COUNT(1) as record_count from user
	INNER JOIN account on user.account_id = account.id where UPPER(account.email) = UPPER(?) AND is_root = true LIMIT 1;`, email)
	
	if recordCount > 0 {
		return 
	}
	
	hashedPassword, hashingErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashingErr != nil {
		logger.Error(hashingErr.Error(), zap.String("error", "hashingErr"))
		return
    }
	transaction, transactErr := db.Begin()
	if transactErr != nil {
		logger.Error(transactErr.Error(), zap.String("error", "transactErr"))
		transaction.Rollback()
		return 
	}
	insertAccountQuery := `INSERT INTO account(email, password, is_root) VALUES(?, ?, ?)`
	result, insertAccountErr := transaction.Exec(insertAccountQuery, email, hashedPassword, true)
	if insertAccountErr != nil {
		logger.Error(insertAccountErr.Error(), zap.String("error", "insertAccountErr"))
		transaction.Rollback()
		return
	}
	accountId, lastInsertedIdErr := result.LastInsertId()
	if lastInsertedIdErr != nil {
		logger.Error(lastInsertedIdErr.Error(), zap.String("error", "lastInsertedIdErr"), zap.Int64("accountId", accountId))
		transaction.Rollback()
		return 
	}
	insertUserQuery := `INSERT INTO user(given_name, middle_name, surname, account_id) VALUES(?, ?, ?, ?)`
	result, insertUserErr := transaction.Exec(insertUserQuery, givenName, middleName, surname, accountId)
	if insertUserErr != nil {
		logger.Error(insertUserErr.Error(), zap.String("error", "insertUserErr"))
		transaction.Rollback()
		return
	}
	transaction.Commit()

}
