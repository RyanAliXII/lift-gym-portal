package model

type User struct {
	Id         int    `json:"id" db:"id"`
	GivenName  string `json:"givenName" db:"given_name"`
	MiddleName string `json:"middleName" db:"middle_name"`
	Surname    string `json:"surname" db:"surname"`
	Email      string `json:"email" db:"email"`
	AccountId   int  `json:"account" db:"account_id"`
	IsRoot 		bool `json:"isRoot,omitempty" db:"is_root"`
	Password   string `json:"password" db:"password"`
	RoleId     int `json:"roleId" db:"role_id"`
	Role 	  string `json:"role" db:"role"`
}
