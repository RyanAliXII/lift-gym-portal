package middlewares

import (
	"net/http"

	"slices"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)


func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc{
	return func (c echo.Context) error {
	
		s , getSessionErr := session.Get("sid", c)
		if getSessionErr != nil {
			s.Options.MaxAge = -1
			s.Save(c.Request(), c.Response())
			return handleResponse(c)
		}
		c.Set("test", "is the same")
		if len(s.Values) ==  0{
			s.Options.MaxAge = -1
			s.Save(c.Request(), c.Response())
			return handleResponse(c)
		}
		return next(c)
	}
}
func handleResponse(c echo.Context) error{
	contentType := c.Request().Header.Get("content-type")
	contentTypeWithJSONResponse := []string{"application/x-www-form-urlencoded,application/json"}
	if(	slices.Contains(contentTypeWithJSONResponse, contentType)){
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code": http.StatusUnauthorized,
			"message": "Unauthorized.",
		})
	}	
	return c.Redirect(http.StatusFound, "/app/login")
}



