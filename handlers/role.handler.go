package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type RoleHandler struct {
}

func (h *RoleHandler) RenderRolePage(c echo.Context) error {

	return c.Render(http.StatusOK, "admin/role/main", Data{
		"csrf" : c.Get("csrf"),
	})
}

func NewRoleHandler () RoleHandler {
	return RoleHandler{
		
	}
}