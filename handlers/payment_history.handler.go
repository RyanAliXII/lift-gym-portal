package handlers

import (
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)



type PaymentHistory struct {
	paymentsHistory repository.PaymentHistory
}
func NewPaymentHistoryHandler () PaymentHistory {
	return  PaymentHistory{
		paymentsHistory: repository.NewPaymentHistory(),
	}
}
func(h * PaymentHistory)RenderClientPaymentHistory(c echo.Context)error{
	contentType := c.Request().Header.Get("content-type")

	if contentType == "application/json"{
		sessionData := c.Get("sessionData")
		session := mysqlsession.SessionData{}
		session.Bind(sessionData)
		payments, err := h.paymentsHistory.GetPaymentHistoryByClient(session.User.Id)
	
		if err != nil{
			logger.Error(err.Error(), zap.String("error", "GetPaymentHistoryByClient"))
			return c.JSON(http.StatusInternalServerError, JSONResponse{
				Status: http.StatusInternalServerError,
				Message: "Unknown error occured.",
			})
		}
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: Data{
				"payments": payments,
			},
			Message: "Payments fetched.",
		})
	}
	return c.Render(http.StatusOK, "client/payment-history/main", Data{})
}


