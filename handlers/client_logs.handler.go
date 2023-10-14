package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ClientLogHandler struct {}
func NewClientLogHandler()ClientLogHandler{
	return ClientLogHandler{}
}
func (h *ClientLogHandler) RenderClientLogPage(c echo.Context) error{
	return c.Render(http.StatusOK, "admin/client-logs/main", Data{
		 "title": "Client Logs",
		 "module": "Client Logs",
	} )
}