package model
type PaymentHistory struct {
	Id int `json:"id"`
	Amount int `json:"amount"`
	Description string `json:"description"`
}