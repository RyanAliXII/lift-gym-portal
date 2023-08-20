package middlewares

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)





func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc{
	return func (c echo.Context) error {
		s , getSessionErr := session.Get("sid", c)
		if getSessionErr != nil {
			return c.Redirect(http.StatusFound, "/login")
		}
		if len(s.Values) ==  0{
			s.Options.MaxAge = -1
			s.Save(c.Request(), c.Response())
			return c.Redirect(http.StatusFound, "/login")
		}
		return next(c)
	}
}



