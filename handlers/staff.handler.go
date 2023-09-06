package handlers

import (
	"lift-fitness-gym/app/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type StaffHandler struct {
}

func (h *StaffHandler) RenderStaffPage(c echo.Context)error {
	return c.Render(http.StatusOK, "admin/staff/main", Data{
		"csrf": c.Get("csrf"),
		"title": "Staffs",
		"module": "Staffs",
	})
}
func (h *StaffHandler)NewStaff (c echo.Context) error {
	staff := model.Staff{}
	err := c.Bind(&staff)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "bindErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	err, fieldErrors := staff.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
			Data: Data{
				"errors": fieldErrors,
			},
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "New staff added.",
	})
}

func NewStaffHandler() StaffHandler{
	return StaffHandler{}
}