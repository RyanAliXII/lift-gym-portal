package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type RegistrationHandler struct {
	clientRepo repository.ClientRepository
}

func (h *RegistrationHandler) RenderRegistrationPage(c echo.Context) error {
	return c.Render(http.StatusOK, "client/registration/main", Data{
		"csrf": c.Get("csrf"),
	})
}

func (h *RegistrationHandler)RegisterClient (c echo.Context) error {
     client := model.Client{}
	 err := c.Bind(&client)
     if err != nil {
		logger.Error(err.Error(), zap.String("error", "bindErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: nil,
			Message: "Unknown error occured.",
		})
	 }
	 err, fields := client.ValidateRegistration()
	 if err != nil {
		logger.Error(err.Error(), zap.String("error", "validationErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": fields,
			},
			Message: "Unknown error occured.",
		})
	 }
	 hashedPassword, err := bcrypt.GenerateFromPassword([]byte(client.Password), bcrypt.DefaultCost)
	 client.Password = string(hashedPassword)

	 err = h.clientRepo.New(client)
	 if err != nil {
		logger.Error(err.Error(), zap.String("error", "NewClientErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	 }
	 return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Success.",
	 })
}
func NewRegistrationHandler() RegistrationHandler {
	return RegistrationHandler{
		clientRepo: repository.NewClientRepository(),
	}
}