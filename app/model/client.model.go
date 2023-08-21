package model

import (
	"fmt"
	"lift-fitness-gym/app/db"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/nyaruka/phonenumbers"
)

type Client struct {
	Id               int    `json:"id" db:"id"`
	GivenName        string `json:"givenName" db:"given_name"`
	MiddleName       string `json:"middleName" db:"middle_name"`
	Surname          string `json:"surname" db:"surname"`
	Email            string `json:"email" db:"email"`
	Password         string `json:"password" db:"password"`
	Address          string `json:"address" db:"address"`
	MobileNumber     string `json:"mobileNumber" db:"mobile_number"`
	DateOfBirth      string `json:"dateOfBirth" db:"date_of_birth"`
	EmergencyContact string `json:"emergencyContact" db:"emergency_contact"`
	Model
}
func (m Client) Validate() (error, map[string]string) {

	db := db.GetConnection()
	return m.Model.ValidationRules(&m, 
		validation.Field(&m.GivenName, validation.Required, validation.Length(1, 255)), 
		validation.Field(&m.MiddleName, validation.Required, validation.Length(1, 255)),
		validation.Field(&m.Surname, validation.Required, validation.Length(1, 255)),
		validation.Field(&m.Address, validation.Required),
		validation.Field(&m.Email, validation.Required, validation.Length(1, 255), is.Email, validation.By(func(value interface{}) error {
			recordCount := 0
			query := `SELECT COUNT(1) as record_count from client
			INNER JOIN account on client.account_id = account.id where UPPER(account.email) = UPPER(?) LIMIT 1;`
			db.Get(&recordCount, query, m.Email)
			if recordCount > 0 {
				return fmt.Errorf("email is already registered")
			}
		
			return nil
		})),	
	    validation.Field(&m.MobileNumber, validation.Required, validation.By(func(value interface{}) error {
			p, _ := phonenumbers.Parse(m.MobileNumber, "PH")
			isValid := phonenumbers.IsValidNumberForRegion(p, "PH")
			if !isValid {
				return fmt.Errorf("invalid number")
			}
			return nil
	 })), 
	validation.Field(&m.EmergencyContact, validation.Required, validation.By(func(value interface{}) error {
		p, _ := phonenumbers.Parse(m.EmergencyContact, "PH")
		isValid := phonenumbers.IsValidNumberForRegion(p, "PH")
		if !isValid {
			return fmt.Errorf("invalid number")
		}
		return nil})))
	
	
}