package model

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
}

func (m Client) Validate() error {
	// return validation.ValidateStruct(&m, validation.Field(&m.GivenName, validation.Required, validation.Length(1, 255)), 
	// validation.Field(&m.l, validation.Required, validation.Length(1, 255),)
	return nil
}