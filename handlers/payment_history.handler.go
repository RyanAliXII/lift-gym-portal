package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)



type PaymentHistory struct {

}
func NewPaymentHistoryHandler () PaymentHistory {
	return  PaymentHistory{}
}
func(h * PaymentHistory)RenderClientPaymentHistory(c echo.Context)error{
	
	return c.Render(http.StatusOK, "client/payment-history/main", Data{})
}


