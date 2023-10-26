package handlers

import (
	"database/sql"
	"lift-fitness-gym/app/pkg/mailer"
	"lift-fitness-gym/app/repository"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)
type PasswordHandler struct {
	userRepo repository.UserRepository
	passwordReset repository.PasswordReset
}
func NewPasswordHandler() PasswordHandler {
	return PasswordHandler{
		userRepo: repository.NewUserRepository(),
		passwordReset: repository.NewPasswordResetRepository(),
	}
}
func (h * PasswordHandler) RenderResetPasswordPage( c echo.Context) error {
	return c.Render(http.StatusOK, "public/password/reset-password", Data{
		"csrf" : c.Get("csrf"),
	})
}
func (h * PasswordHandler) RenderResetClientPasswordPage( c echo.Context) error {
	return c.Render(http.StatusOK, "public/password/client/reset-password", Data{
		"csrf" : c.Get("csrf"),
	})
}
func (h * PasswordHandler)RenderChangePasswordPage(c echo.Context) error {
	key := c.QueryParam("key")

	_, err := h.passwordReset.GetByPublicKey(key)
	if err != nil {
		return c.Render(http.StatusNotFound, "partials/error/404-page", nil)
	}
	return c.Render(http.StatusOK,"public/password/change-password", Data{
		"csrf" : c.Get("csrf"),
	})
}
func (h * PasswordHandler) ResetPassword( c echo.Context) error {
	email := c.FormValue("email")
	err := validation.Validate(&email, validation.Required.Error("Email is required."), is.Email.Error("Email is not valid."))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "validation error"))
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
		    Message: "OK",
		})
	}
	user, err := h.userRepo.GetUserByEmail(email)
	
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "GetUserByEmail"))
		if err != sql.ErrNoRows {
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


func (h * PasswordHandler) ResetClientPassword( c echo.Context) error {
	email := c.FormValue("email")
	err := validation.Validate(&email, validation.Required.Error("Email is required."), is.Email.Error("Email is not valid."))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "validation error"))
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
		    Message: "OK",
		})
	}
	user, err := h.userRepo.GetClientUserByEmail(email)
	
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "GetUserByEmail"))
		if err != sql.ErrNoRows {
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

func (h * PasswordHandler) ChangePassword( c echo.Context) error {
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirmPassword")
	publicKey := c.FormValue("publicKey")
	err := validation.Validate(&password, validation.Required.Error("Password is required."), validation.Length(10, 30).Error("Password length must be atleast 10 characters to 30 characters."))
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{Status: http.StatusBadRequest, 
			Data: Data{
			"errors":Data{
					"password" : err.Error(),			
				},
			},
			Message: "Validation error."})	
	}
	if password != confirmPassword {
		return c.JSON(http.StatusBadRequest, JSONResponse{Status: http.StatusBadRequest, 
			
			Data: Data{
				"errors": Data{
					"confirmPassword": "Password and confirm password is not the same.",
				},
			},
			Message: "Validation error."})
	}
	hashedPWD, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "hashingErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	err = h.passwordReset.ChangePasswordByPublicKey(string(hashedPWD), publicKey)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "passwordResetErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{Status:http.StatusOK, Message: "OK" })
}
