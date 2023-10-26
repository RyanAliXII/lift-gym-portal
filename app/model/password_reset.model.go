package model

type PasswordReset struct {
	Id string `json:"id" db:"id"`
	PublicId string `json:"publicId" db:"public_id"`
	AccountId int `json:"accountId" db:"account_id"`
}