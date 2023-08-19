package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type LoginHandler struct {
}


type LoginHandlerInterface interface {
	RenderLoginPage(c echo.Context ) error 

}
func (h *LoginHandler) RenderLoginPage(c echo.Context) error{
	return c.Render(http.StatusOK, "login/main", Data{
		"title": "Sign In",
	} )
}

func NewLoginHandler() LoginHandlerInterface {
	return &LoginHandler{}
}




