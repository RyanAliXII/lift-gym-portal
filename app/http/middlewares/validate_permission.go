package middlewares

import (
	"lift-fitness-gym/app/pkg/acl"
	"lift-fitness-gym/app/pkg/mysqlsession"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ValidatePermissions(requiredPermission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session := mysqlsession.SessionData{}
			session.Bind(c.Get("sessionData"))
			hasPermission := acl.HasPermission(requiredPermission, session.User.Permissions)
			if hasPermission {
				return next(c)
			}
			contentType := c.Request().Header.Get("Content-Type")
			if contentType == "application/json"{
				return c.JSON(http.StatusForbidden, "You are not allow toa access this page.")
			}
			return c.Render(http.StatusForbidden, "partials/error/403-page", nil)
		}
	}
}