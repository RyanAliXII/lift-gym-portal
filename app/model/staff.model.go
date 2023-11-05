package model

import (
	"fmt"
	"lift-fitness-gym/app/db"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Staff struct {
	Id         int    `json:"id" db:"id"`
	GivenName  string `json:"givenName" db:"given_name"`
	MiddleName string `json:"middleName" db:"middle_name"`
	Surname    string `json:"surname" db:"surname"`
	AccountId  string `json:"accountId,omitempty" db:"account_id"`
	Email      string `json:"email" db:"email"`
	Password   string `json:"password,omitempty" db:"password"`
	Address          string `json:"address,omitempty" db:"address"`
	MobileNumber     string `json:"mobileNumber,omitempty" db:"mobile_number"`
	PublicId 		string `json:"publicId" db:"public_id"`
	Gender 			string `json:"gender" db:"gender"`
	EmergencyContact string `json:"emergencyContact,omitempty" db:"emergency_contact"`
	DateOfBirth      string `json:"dateOfBirth,omitempty" db:"date_of_birth"`
	RoleId	   int `json:"roleId" db:"role_id"`
	Model
}
func (m Staff) Validate() (error, map[string]string) {
	db := db.GetConnection()
	return m.Model.ValidationRules(&m, 
		validation.Field(&m.GivenName, validation.Required.Error("Given name is required."), validation.Length(1, 255).Error("Given name must be 1 to 255 characters.")),
		validation.Field(&m.Surname, validation.Required.Error("Surname is required."), validation.Length(1, 255).Error("Surname must be 1 to 255 characters.")),
		validation.Field(&m.MiddleName, validation.Required.Error("Middle name is required."), validation.Length(1, 255).Error("Middle name must be 1 to 255 characters.")),
		validation.Field(&m.RoleId, validation.Required.Error("Role is required."), validation.Min(1).Error("Role is required.")),
		validation.Field(&m.Email, validation.Required.Error("Email is required."), validation.Length(1, 255).Error("Email must be 1 to 255 characters"), is.Email.Error("Invalid email"), validation.By(func(value interface{}) error {
			recordCount := 0
			query := `SELECT COUNT(1) as record_count from user
			INNER JOIN account on user.account_id = account.id where UPPER(account.email) = UPPER(?) LIMIT 1;`
			db.Get(&recordCount, query, m.Email)
			if recordCount > 0 {
				return fmt.Errorf("Email is already registered.")
			}		
			return nil
		})),
	)
}
func (m Staff) ValidateUpdate() (error, map[string]string) {
	db := db.GetConnection()
	return m.Model.ValidationRules(&m, 
		validation.Field(&m.GivenName, validation.Required.Error("Given name is required."), validation.Length(1, 255).Error("Given name must be 1 to 255 characters.")),
		validation.Field(&m.Surname, validation.Required.Error("Surname is required."), validation.Length(1, 255).Error("Surname must be 1 to 255 characters.")),
		validation.Field(&m.MiddleName, validation.Required.Error("Middle name is required."), validation.Length(1, 255).Error("Middle name must be 1 to 255 characters.")),
		validation.Field(&m.RoleId, validation.Required.Error("Role is required."), validation.Min(1).Error("Role is required.")),
		validation.Field(&m.Email, validation.Required.Error("Email is required."), validation.Length(1, 255).Error("Email must be 1 to 255 characters"), is.Email.Error("Invalid email"), validation.By(func(value interface{}) error {
			recordCount := 0
			query := `SELECT COUNT(1) as record_count from user
			INNER JOIN account on user.account_id = account.id where UPPER(account.email) = UPPER(?) AND user.id != ? LIMIT 1;`
			db.Get(&recordCount, query, m.Email, m.Id)
			if recordCount > 0 {
				return fmt.Errorf("Email is already registered.")
			}		
			return nil
		})),
	)
}

