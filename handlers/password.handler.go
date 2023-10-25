package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)
type PasswordHandler struct {
	
}
func (h * PasswordHandler) RenderResetPasswordPage( c echo.Context) error {
	

	return c.Render(http.StatusOK, "public/password/reset-password",Data{
		"csrf" : c.Get("csrf"),
	})
}
