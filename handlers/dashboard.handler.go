package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {


}
func (h *DashboardHandler) RenderDashboardPage(c echo.Context) error{

	return c.Render(http.StatusOK, "dashboard/main", Data{
		"title": "Dashboard",
		"module": "Dashboard",
	} )
}

func NewDashboardHandler() DashboardHandler{
	return DashboardHandler{}
}




