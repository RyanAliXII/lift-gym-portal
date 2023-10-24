package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"lift-fitness-gym/app/db"
	"time"

	_ "time/tzdata"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/nyaruka/phonenumbers"
)

type Coach struct {
	Id               int    `json:"id" db:"id"`
	GivenName        string `json:"givenName" db:"given_name"`
	MiddleName       string `json:"middleName" db:"middle_name"`
	Surname          string `json:"surname" db:"surname"`
	Email            string `json:"email" db:"email"`
	Password         string `json:"password,omitempty" db:"password"`
	Address          string `json:"address" db:"address"`
	Description 	 string `json:"description" db:"description"`
	MobileNumber     string `json:"mobileNumber" db:"mobile_number"`
	DateOfBirth      string `json:"dateOfBirth" db:"date_of_birth"`
	AccountId 		 int `json:"accountId" db:"account_id"`
	IsVerified		 bool `json:"isVerified" db:"is_verified"`
	EmergencyContact string `json:"emergencyContact" db:"emergency_contact"`
	Images 			 CoachImages `json:"images" db:"images"`
	Model
}
func (m Coach) Validate() (error, map[string]string) {
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
			query := `SELECT COUNT(1) as record_count from coach
			INNER JOIN account on coach.account_id = account.id where UPPER(account.email) = UPPER(?) LIMIT 1;`
			db.Get(&recordCount, query, m.Email)
			if recordCount > 0 {
				return fmt.Errorf("email is already registered")
			}
			return nil
		})),	
	    validation.Field(&m.MobileNumber, validation.Required.Error("Mobile number is required."), validation.By(func(value interface{}) error {
			p, _ := phonenumbers.Parse(m.MobileNumber, "PH")
			isValid := phonenumbers.IsValidNumberForRegion(p, "PH")
			if !isValid {
				return fmt.Errorf("invalid number")
			}
			return nil
		 })), 
		validation.Field(&m.EmergencyContact, validation.Required.Error("Emergency contact number is required."), validation.By(func(value interface{}) error {
			p, _ := phonenumbers.Parse(m.EmergencyContact, "PH")
			isValid := phonenumbers.IsValidNumberForRegion(p, "PH")
			if !isValid {
				return fmt.Errorf("invalid number")
			}
			return nil})))
}
func (m Coach) ValidateUpdate() (error, map[string]string) {
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
		validation.Field(&m.Email, validation.Required.Error("Email is required."), validation.Length(1, 255).Error("Email is required."), is.Email.Error("Invalid email."), validation.By(func(value interface{}) error {
			recordCount := 0
			query := `SELECT COUNT(1) as record_count from coach
			INNER JOIN account on coach.account_id = account.id where UPPER(account.email) = UPPER(?) and coach.id != ? LIMIT 1;`
			db.Get(&recordCount, query, m.Email, m.Id)
			if recordCount > 0 {
				return fmt.Errorf("email is already registered")
			}		
			return nil
		})),	
		validation.Field(&m.MobileNumber, validation.Required.Error("Mobile number is required."), validation.By(func(value interface{}) error {
			p, _ := phonenumbers.Parse(m.MobileNumber, "PH")
			isValid := phonenumbers.IsValidNumberForRegion(p, "PH")
			if !isValid {
				return fmt.Errorf("invalid number")
			}
			return nil
		 })), 
		validation.Field(&m.EmergencyContact, validation.Required.Error("Emergency contact number is required."), validation.By(func(value interface{}) error {
			p, _ := phonenumbers.Parse(m.EmergencyContact, "PH")
			isValid := phonenumbers.IsValidNumberForRegion(p, "PH")
			if !isValid {
				return fmt.Errorf("invalid number")
			}
			return nil})))
}


type CoachJSON struct {
	Coach
}
func (instance *CoachJSON) Scan(value interface{}) error {
	val, valid := value.([]byte)
	if valid {
		unmarshalErr := json.Unmarshal(val, instance)
		if unmarshalErr != nil {
			*instance = CoachJSON{}
		}
	} else {
		*instance = CoachJSON{}
	}
	return nil

}
func (copy CoachJSON) Value(value interface{}) (driver.Value, error) {
	return copy, nil
}

type HiredCoach struct {
	Id int `json:"id" db:"id"`
	CoachId int `json:"coachId" db:"coach_id"`
	RateId int `json:"rateId" db:"rate_id"`
	RateSnapshotId int `json:"rateSnapshotId" db:"rate_snapshot_id"`
	ClientId int `json:"clientId" db:"client_id"`
	Coach CoachJSON `json:"coach" db:"coach"`
	Client ClientJSON `json:"client" db:"client"`
	Rate CoachRateJSON `json:"rate" db:"rate"`
	RateSnapshot CoachRateJSON `json:"rateSnapshot" db:"rate_snapshot"`
	Status string `json:"status" db:"status"`
	StatusId int `json:"statusId" db:"status_id"`
	Remarks string `json:"remarks" db:"remarks"`
	MeetingTime string `json:"meetingTime" db:"meeting_time"`
	Model
}


func(m HiredCoach) Validate() (error, map[string]string) {
	return m.Model.ValidationRules(&m,
		 validation.Field(&m.CoachId, validation.Required.Error("Coach is required."), validation.Min(1).Error("Coach is required.")), 
		 validation.Field(&m.RateId, validation.Required.Error("Rate is required."), validation.Min(1).Error("Rate is required.")))
}

func(m HiredCoach) ValidateMeetingTime() (error, map[string]string) {
	return m.Model.ValidationRules(&m,
		 validation.Field(&m.MeetingTime, validation.Required.Error("Datetime is required."), validation.By(func(value interface{}) error {
			t, err := time.Parse(time.RFC3339, value.(string))
			if err != nil {
				return fmt.Errorf("Meeting time is required.")
			}
			location, err := time.LoadLocation("Asia/Manila")
			if err != nil {
				return fmt.Errorf("Unknown error occured")
			}
			now := time.Now().In(location)
			t = t.In(location)

			if t.Before(now) {
				return fmt.Errorf("Date cannot be past current date.")
			}
			return nil
		 })))
}

