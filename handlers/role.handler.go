package handlers

import (
	"fmt"
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/acl"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RoleHandler struct {
}

func (h *RoleHandler) RenderRolePage(c echo.Context) error {

	return c.Render(http.StatusOK, "admin/role/main", Data{
		"csrf" : c.Get("csrf"),
		"permissions": acl.Permissions,
	})
}

func (h *RoleHandler) NewRole(c echo.Context) error {
	role := model.Role{}

	err := c.Bind(&role)
	if err != nil {

		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	err, fields := role.Validate()

	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": fields,
			},
			Message: "Validation error.",
		})
	}
	fmt.Println(role)
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Success",
	})
}

func NewRoleHandler () RoleHandler {
	return RoleHandler{}
}