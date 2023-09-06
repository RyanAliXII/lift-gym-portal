package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
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

	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "New staff added.",
	})
}

func NewStaffHandler() StaffHandler{
	return StaffHandler{}
}