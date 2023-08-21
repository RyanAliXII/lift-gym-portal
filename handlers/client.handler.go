package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ClientHandler struct {
}
func (h * ClientHandler) RenderClientPage(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/clients/main", Data{
		"title": "Clients",
		"module": "Clients",
	})
}
func (h * ClientHandler) RenderClientRegistrationForm(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/clients/register-client", Data{
		"title": "Client | Registration ",
		"module": "Registration Form",
	})
}


func NewClientHandler() ClientHandler {
	return ClientHandler{}
}