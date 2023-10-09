package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {


}
func (h *DashboardHandler) RenderDashboardPage(c echo.Context) error{
	return c.Render(http.StatusOK, "admin/dashboard/main", Data{
		"title": "Dashboard",
		"module": "Dashboard",
	} )
}
func (h *DashboardHandler) RenderClientDashboard(c echo.Context ) error {
	return c.Render(http.StatusOK, "client/dashboard/main", Data{
		"title": "Dashboard",
		"module": "Dashboard",
	} )

}
func (h * DashboardHandler) RenderCoachDashboard(c echo.Context) error {
	return c.Render(http.StatusOK, "coach/dashboard/main", Data{
		"title": "Dashboard",
		"module": "Dashboard",
} )
}
func NewDashboardHandler() DashboardHandler{
	return DashboardHandler{}
}




