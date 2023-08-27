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

}
func(h * VerificationHandler) VerifyEmail(c echo.Context) error {
	id := c.Param("id")
	verification, err := h.verificationRepo.GetEmailVerificationByPublicId(id)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "getEmailVerificationByPublicId"))
		return c.JSON(http.StatusNotFound, JSONResponse{
			Status: http.StatusNotFound,
			Message: "Page not found.",
		})
	}
	fmt.Println(verification)
	return nil
}


func NewVerificationHandler() VerificationHandler{

	return VerificationHandler{
		verificationRepo: repository.NewVerificationRepository(),
	}
}