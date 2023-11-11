package repository

import "github.com/jmoiron/sqlx"


type PaymentHistory struct {
	db  sqlx.DB
}
func(repo  * PaymentHistory)GetPaymentHistoryByClient(clientId int )error {

	return nil
}