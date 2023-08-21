package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ClientHandler struct {
}
func (h * ClientHandler) RenderMemberHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/clients/main", Data{
		"title": "Clients",
		"module": "Clients",
	})
}


func NewClientHandler() ClientHandler {
	return ClientHandler{}
}