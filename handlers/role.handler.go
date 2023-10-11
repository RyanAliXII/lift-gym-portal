package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/acl"
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type RoleHandler struct {
	roleRepo repository.RoleRepository
}

func (h *RoleHandler) RenderRolePage(c echo.Context) error {
	contentType := c.Request().Header.Get("Content-Type")
	if contentType == "application/json"{
		roles, err := h.roleRepo.GetRoles()
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "getRoles"))
		}
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: Data{
				"roles": roles,
			},
			Message: "Roles fetched.",
		})
	}	
	return c.Render(http.StatusOK, "admin/role/main", Data{
		"csrf" : c.Get("csrf"),
		"title": "Access Control | Roles",
		"module":"Roles and Permissions",
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
	err = h.roleRepo.NewRole(role)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "newRoleErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Success",
	})
}

func NewRoleHandler () RoleHandler {
	return RoleHandler{
		roleRepo: repository.NewRoleRepository(),
	}
}