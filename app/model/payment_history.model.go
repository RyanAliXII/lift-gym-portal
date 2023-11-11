package model
type PaymentHistory struct {
	Amount float64 `json:"amount" db:"amount"`
	ClientId int `json:"client_id" db:"client_id"`
	Description string `json:"description" db:"description"`
	CreatedAt string `json:"createdAt" db:"created_at"`
}