package handlers

import (
	"fmt"
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)



type VerificationHandler struct {

	verificationRepo repository.VerificationRepository
	clientRepo repository.ClientRepository

}
func(h * VerificationHandler) VerifyEmail(c echo.Context) error {
	id := c.Param("id")
	verification, err := h.verificationRepo.GetEmailVerificationByPublicId(id)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "GetEmailVerificationByPublicId"))
		return c.Render(http.StatusNotFound, "partials/error/404-page", nil)
		// return c.JSON(http.StatusNotFound, JSONResponse{
		// 	Status: http.StatusNotFound,
		// 	Message: "Page not found.",
		// })
	}
	fmt.Println(verification)
	err = h.verificationRepo.MarkAsComplete(verification.Id) 
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "MarkAsComplete"))
		return c.Render(http.StatusNotFound, "partials/error/404-page", nil)
	}
	err = h.clientRepo.MarkAsVerified(verification.ClientId)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "MarkAsVerified"))
		return c.Render(http.StatusNotFound, "partials/error/404-page", nil)
	}
	return c.Render(http.StatusOK, "verification/main", nil)
}


func NewVerificationHandler() VerificationHandler{
	return VerificationHandler{
		verificationRepo: repository.NewVerificationRepository(),
		clientRepo: repository.NewClientRepository(),
	}
}