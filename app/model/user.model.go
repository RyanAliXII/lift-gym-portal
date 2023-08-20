package model

type User struct {
	Id         int    `json:"id" db:"id"`
	GivenName  string `json:"givenName" db:"given_name"`
	MiddleName string `json:"middleName" db:"middle_name"`
	Surname    string `json:"surname" db:"surname"`
	Email      string `json:"email" db:"email"`
	Password   string `json:"password" db:"password"`
}