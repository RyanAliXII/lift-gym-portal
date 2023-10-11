package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type RegistrationHandler struct {
}

func (h *RegistrationHandler) RenderRegistrationPage(c echo.Context) error {
	return c.Render(http.StatusOK, "client/dashboard/main", Data{})
}


func NewRegistrationHandler() RegistrationHandler {
	return RegistrationHandler{}
}