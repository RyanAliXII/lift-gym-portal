package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"lift-fitness-gym/app/db"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/nyaruka/phonenumbers"
)

type Client struct {
	Id               int    `json:"id" db:"id"`
	GivenName        string `json:"givenName" db:"given_name"`
	MiddleName       string `json:"middleName" db:"middle_name"`
	Surname          string `json:"surname" db:"surname"`
	Email            string `json:"email,omitempty" db:"email"`
	Password         string `json:"password,omitempty" db:"password"`
	Address          string `json:"address,omitempty" db:"address"`
	MobileNumber     string `json:"mobileNumber,omitempty" db:"mobile_number"`
	DateOfBirth      string `json:"dateOfBirth,omitempty" db:"date_of_birth"`
	AccountId 		 string `json:"accountId,omitempty" db:"account_id"`
	IsVerified		 bool 	`json:"isVerified" db:"is_verified"`
	IsMember		 bool 	`json:"isMember" db:"is_member"`
	EmergencyContact string `json:"emergencyContact,omitempty" db:"emergency_contact"`
	Model
}

func (m Client) Validate() (error, map[string]string) {

	db := db.GetConnection()
	return m.Model.ValidationRules(&m, 
		validation.Field(&m.GivenName, validation.Required.Error("Given name is required."), validation.Length(1, 255).Error("Given name is required.")), 
		validation.Field(&m.MiddleName, validation.Required.Error("Middle name is required."), validation.Length(1, 255).Error("Middle name is required.")), 
		validation.Field(&m.Surname, validation.Required.Error("Surname is required."), validation.Length(1, 255).Error("Surname is required.")), 
		validation.Field(&m.DateOfBirth, validation.Required.Error("Date of birth is required."), validation.By(func(value interface{}) error {
			format := "2006-01-02"
			strDate ,_ := value.(string)
			_, err := time.Parse(format, strDate)
			if err != nil {
				return fmt.Errorf("Date of birth is required.")
			}
			return nil
		})),
		validation.Field(&m.Address, validation.Required.Error("Address is required."), validation.Length(1, 255).Error("Address should be atleast 1 to 255 characters long")),
		validation.Field(&m.Email, validation.Required.Error("Email is required."), validation.Length(1, 255).Error("Email is required."), is.Email.Error("Invalid email"), validation.By(func(value interface{}) error {
			recordCount := 0
			query := `SELECT COUNT(1) as record_count from client
			INNER JOIN account on client.account_id = account.id where UPPER(account.email) = UPPER(?) and deleted_at is null LIMIT 1;`
			db.Get(&recordCount, query, m.Email)
			if recordCount > 0 {
				return fmt.Errorf("Email is already registered")
			}
		
			return nil
		})),	
	    validation.Field(&m.MobileNumber, validation.Required.Error("Mobile number is required."), validation.By(func(value interface{}) error {
			p, _ := phonenumbers.Parse(m.MobileNumber, "PH")
			isValid := phonenumbers.IsValidNumberForRegion(p, "PH")
			if !isValid {
				return fmt.Errorf("Invalid number")
			}
			return nil
		 })),
 
		validation.Field(&m.EmergencyContact, validation.Required.Error("Emergency contact number is required."), validation.By(func(value interface{}) error {
			p, _ := phonenumbers.Parse(m.EmergencyContact, "PH")
			isValid := phonenumbers.IsValidNumberForRegion(p, "PH")
			if !isValid {
				return fmt.Errorf("Invalid number")
			}
			return nil})))
}
func (m Client) ValidateUpdate() (error, map[string]string) {

	db := db.GetConnection()
	return m.Model.ValidationRules(&m, 
		validation.Field(&m.GivenName, validation.Required.Error("Given name is required."), validation.Length(1, 255).Error("Given name is required.")), 
		validation.Field(&m.MiddleName, validation.Required.Error("Middle name is required."), validation.Length(1, 255).Error("Middle name is required.")), 
		validation.Field(&m.Surname, validation.Required.Error("Surname is required."), validation.Length(1, 255).Error("Surname is required.")), 
		validation.Field(&m.DateOfBirth, validation.Required.Error("Date of birth is required."), validation.By(func(value interface{}) error {
			format := "2006-01-02"
			strDate ,_ := value.(string)
			_, err := time.Parse(format, strDate)
			if err != nil {
				return fmt.Errorf("Date of birth is required.")
			}
			return nil
		})),
		validation.Field(&m.Address, validation.Required.Error("Address is required."), validation.Length(1, 255).Error("Address should be atleast 1 to 255 characters long.")),
		validation.Field(&m.Email, validation.Required.Error("Email is required."), validation.Length(1, 255).Error("Email is required."), is.Email.Error("Invalid email"), validation.By(func(value interface{}) error {
			recordCount := 0
			query := `SELECT COUNT(1) as record_count from client
			INNER JOIN account on client.account_id = account.id where UPPER(account.email) = UPPER(?) and client.id != ?  and deleted_at is null LIMIT 1;`
			db.Get(&recordCount, query, m.Email, m.Id)
			if recordCount > 0 {
				return fmt.Errorf("Email is already registered.")
			}		
			return nil
		})),	
	    validation.Field(&m.MobileNumber, validation.Required.Error("Mobile number is required."), validation.By(func(value interface{}) error {
			p, _ := phonenumbers.Parse(m.MobileNumber, "PH")
			isValid := phonenumbers.IsValidNumberForRegion(p, "PH")
			if !isValid {
				return fmt.Errorf("Invalid number")
			}
			return nil
		 })), 
		validation.Field(&m.EmergencyContact, validation.Required.Error("Emergency contact number is required."), validation.By(func(value interface{}) error {
			p, _ := phonenumbers.Parse(m.EmergencyContact, "PH")
			isValid := phonenumbers.IsValidNumberForRegion(p, "PH")
			if !isValid {
				return fmt.Errorf("Invalid number")
			}
			return nil})))
}
func (m Client) ValidateRegistration () (error, map[string]string){

	db := db.GetConnection()
	return m.Model.ValidationRules(&m, 
		validation.Field(&m.GivenName, validation.Required.Error("Given name is required."), validation.Length(1, 255).Error("Given name is required.")), 
		validation.Field(&m.MiddleName, validation.Required.Error("Middle name is required."), validation.Length(1, 255).Error("Middle name is required.")), 
		validation.Field(&m.Surname, validation.Required.Error("Surname is required."), validation.Length(1, 255).Error("Surname is required.")), 
		validation.Field(&m.Password, validation.Required.Error("Password is required."), validation.Length(10, 30).Error("Password must be atleast 10 characters to 30 characters long.")),
		validation.Field(&m.DateOfBirth, validation.Required.Error("Date of birth is required."), validation.By(func(value interface{}) error {
			format := "2006-01-02"
			strDate ,_ := value.(string)
			_, err := time.Parse(format, strDate)
			if err != nil {
				return fmt.Errorf("Date of birth is required.")
			}
			return nil
		})),
	validation.Field(&m.Email, validation.Required.Error("Email is required."), validation.Length(1, 255).Error("Email is required."), is.Email.Error("Invalid email"), validation.By(func(value interface{}) error {
			recordCount := 0
			query := `SELECT COUNT(1) as record_count from client
			INNER JOIN account on client.account_id = account.id where UPPER(account.email) = UPPER(?) and deleted_at is null LIMIT 1;`
			db.Get(&recordCount, query, m.Email)
			if recordCount > 0 {
				return fmt.Errorf("Email is already registered")
			}
		
			return nil
		})),	
	)
}
type ClientJSON struct {
	Client
}

func (instance *ClientJSON) Scan(value interface{}) error {
	val, valid := value.([]byte)
	if valid {
		unmarshalErr := json.Unmarshal(val, instance)
		if unmarshalErr != nil {
			*instance = ClientJSON{}
		}
	} else {
		*instance = ClientJSON{}
	}
	return nil

}
func (copy ClientJSON) Value(value interface{}) (driver.Value, error) {
	return copy, nil
}