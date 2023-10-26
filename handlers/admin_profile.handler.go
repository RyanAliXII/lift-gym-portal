package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)



type AdminProfileHandler struct {

}

func NewAdminProfileHandler() AdminProfileHandler {
	return AdminProfileHandler{}
}
func (h * AdminProfileHandler) RenderAdminProfile (c echo.Context) error {
	return c.Render(http.StatusBadRequest, "admin/profile/main",nil)
}