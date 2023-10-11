package handlers

import (
	"lift-fitness-gym/app/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type RegistrationHandler struct {
}

func (h *RegistrationHandler) RenderRegistrationPage(c echo.Context) error {
	return c.Render(http.StatusOK, "client/registration/main", Data{})
}

func (h *RegistrationHandler)registerClient (c echo.Context) error {
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
	 return nil
}
func NewRegistrationHandler() RegistrationHandler {
	return RegistrationHandler{}
}