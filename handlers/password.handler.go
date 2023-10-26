package handlers

import (
	"database/sql"
	"lift-fitness-gym/app/pkg/mailer"
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)
type PasswordHandler struct {
	userRepo repository.UserRepository
	passwordReset repository.PasswordReset
}
func NewPasswordHandler() PasswordHandler {
	return PasswordHandler{
		userRepo: repository.NewUserRepository(),
	}
}
func (h * PasswordHandler) RenderResetPasswordPage( c echo.Context) error {
	return c.Render(http.StatusOK, "public/password/reset-password", Data{
		"csrf" : c.Get("csrf"),
	})
}
func (h * PasswordHandler) ResetPassword( c echo.Context) error {
	email := c.FormValue("email")
	user, err := h.userRepo.GetUserByEmail(email)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Error(err.Error(), zap.String("error", "GetUserByEmail"))
			return c.JSON(http.StatusInternalServerError,  JSONResponse{
				Status: http.StatusInternalServerError,
				Message: "Unknown error occured.",
			})
		}
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
		    Message: "OK",
		})
	}
	passwordReset, err := h.passwordReset.New(user.AccountId)
	if err != nil {
			logger.Error(err.Error(), zap.String("error", "NewPasswordReset"))
	}
	go mailer.SendEmailPasswordReset([]string{user.Email}, user.GivenName,  passwordReset.PublicId)
	return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
		    Message: "OK",
	})
}

