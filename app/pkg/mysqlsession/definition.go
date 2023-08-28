package mysqlsession

import (
	"encoding/json"
	"errors"
)

type SessionData struct {
	User SessionUser `json:"user"`
}


type SessionUser struct {
	Id         int    `json:"id" db:"id"`
	GivenName  string `json:"givenName" db:"given_name"`
	MiddleName string `json:"middleName" db:"middle_name"`
	Surname    string `json:"surname" db:"surname"`
	Email      string `json:"email" db:"email"`
}
func(s SessionData) ToBytes() ([]byte, error) {
	sessionDataBytes,err := json.Marshal(s)
	return  sessionDataBytes, err
}
func(s * SessionData )Bind(v interface{}) error {
	bytes, isByte := v.([]byte)
	if !isByte {
		return errors.New("invalid value passed. expecting a bytes value.")
	}
	err := json.Unmarshal(bytes, s)
	return err
}


