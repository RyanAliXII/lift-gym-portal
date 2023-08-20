package db

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func CreateRootAccount() {
	givenName := os.Getenv("ROOT_USER_GIVEN_NAME")
	middleName := os.Getenv("ROOT_USER_MIDDLE_NAME")
	surname := os.Getenv("ROOT_USER_SURNAME")
	email := os.Getenv("ROOT_USER_EMAIL")
	password := os.Getenv("ROOT_USER_PASSWORD")
	if len(givenName) == 0 || len(middleName) == 0 || len(surname) == 0 || len(email) == 0 || len(password) == 0{
		fmt.Println("One or more fields are empty.")
		return
	}
	hashedPassword, hashingErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashingErr != nil {
        fmt.Println(hashingErr)
		return
    }
	transaction, transactErr := db.Begin()
	if transactErr != nil {
		transaction.Rollback()
		return 
	}
	insertAccountQuery := `INSERT INTO account(email, password) VALUES(?, ?)`
	result, insertAccountErr := transaction.Exec(insertAccountQuery, email, hashedPassword)
	if insertAccountErr != nil {
		fmt.Println(insertAccountErr.Error())
		transaction.Rollback()
		return
	}
	accountId, lastInsertedIdErr := result.LastInsertId()
	if lastInsertedIdErr != nil {
		fmt.Println(lastInsertedIdErr.Error())
		transaction.Rollback()
		return 
	}
	insertUserQuery := `INSERT INTO user(given_name, middle_name, surname, account_id) VALUES(?, ?, ?, ?)`
	result, insertUserErr := transaction.Exec(insertUserQuery, givenName, middleName, surname, accountId)
	if insertUserErr != nil {
		fmt.Println(insertUserErr.Error())
		transaction.Rollback()
		return
	}
	transaction.Commit()

}
