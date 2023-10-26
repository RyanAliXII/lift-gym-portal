package handlers

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)


type LogoutHandler struct {

}

func NewLogoutHandler() LogoutHandler{
	return LogoutHandler{}
}
func (h * LogoutHandler)LogoutAdmin (c echo.Context) error {
	s , err := session.Get("sid", c)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "getSessionErr"))
	}
	s.Options.MaxAge = -1
	s.Values["data"] = nil
	s.Save(c.Request(),c.Response())
	return c.JSON(http.StatusOK, "OK")
}

func (h * LogoutHandler)LogoutClient (c echo.Context) error {
	s , err := session.Get("client_sid", c)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "getSessionErr"))
	}
	s.Options.MaxAge = -1
	s.Values["data"] = nil
	s.Save(c.Request(),c.Response())
	return c.JSON(http.StatusOK, "OK")
}
func (h * LogoutHandler)LogoutCoach (c echo.Context) error {
	s , err := session.Get("coach_sid", c)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "getSessionErr"))
	}
	s.Options.MaxAge = -1
	s.Values["data"] = nil
	s.Save(c.Request(),c.Response())
	return c.JSON(http.StatusOK, "OK")
}