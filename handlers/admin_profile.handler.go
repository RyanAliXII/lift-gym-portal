package handlers

import (
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/app/repository"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)



type AdminProfileHandler struct {
	userRepo repository.UserRepository
}

func NewAdminProfileHandler() AdminProfileHandler {
	return AdminProfileHandler{
		userRepo: repository.NewUserRepository(),
	}
}
func (h * AdminProfileHandler) RenderAdminProfile (c echo.Context) error {
	return c.Render(http.StatusBadRequest, "admin/profile/main",nil)
}

func (h * AdminProfileHandler) ChangePassword (c echo.Context) error {
	oldPassword := c.FormValue("oldPassword")
	newPassword := c.FormValue("newPassword")
	sessionData := mysqlsession.SessionData{}
	sessionData.Bind(c.Get("sessionData"))
	user, err := h.userRepo.GetUserById(sessionData.User.Id)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "GetUserByIdErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{Status: http.StatusBadRequest, 
			Data: Data{
			"errors":Data{
					"oldPassword" : "Old password is incorrect.",			
				},
			},
			Message: "Validation error."})	
	}


	err = validation.Validate(&newPassword, validation.Required.Error("New Password is required."), validation.Length(10, 30).Error("New password length must be atleast 10 characters to 30 characters."))
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{Status: http.StatusBadRequest, 
			Data: Data{
			"errors":Data{
					"newPassword" : err.Error(),			
				},
			},
			Message: "Validation error."})	
	}

	
	
	hashedPWD, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "hashingErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	err = h.userRepo.UpdateAdminPasswordByAccountId(user.AccountId, string(hashedPWD))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "updatePWDErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "OK",
	})
}