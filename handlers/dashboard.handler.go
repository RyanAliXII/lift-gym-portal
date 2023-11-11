package handlers

import (
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type DashboardHandler struct {
	dashboardRepo repository.Dashboard
}
func (h *DashboardHandler) RenderDashboardPage(c echo.Context) error{
	contentType := c.Request().Header.Get("content-type")
	if contentType == "application/json"{
		data, err := h.dashboardRepo.GetAdminDashboardData()
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "GetAdminDashboardData"))
		}
		data.MonthlyWalkIns, err = h.dashboardRepo.GetMonthlyWalkIns()
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "GetMonthlyWalkIns"))
		}
		data.WeeklyWalkIns, err = h.dashboardRepo.GetWeeklyWalkIns()
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "GetWeeklyWalkIns"))
		}
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: Data{
				"dashboardData": data,

			},
			Message: "Dashboard data fetched.",
		})
	}
	return c.Render(http.StatusOK, "admin/dashboard/main", Data{
		"title": "Dashboard",
		"module": "Dashboard",
		"csrf": c.Get("csrf"),
	} )
}
func (h *DashboardHandler) RenderClientDashboard(c echo.Context ) error {
	return c.Render(http.StatusOK, "client/dashboard/main", Data{
		"title": "Dashboard",
		"module": "Dashboard",
		"csrf": c.Get("csrf"),
	} )

}
func (h * DashboardHandler) RenderCoachDashboard(c echo.Context) error {
	return c.Render(http.StatusOK, "coach/dashboard/main", Data{
		"title": "Dashboard",
		"module": "Dashboard",
		"csrf": c.Get("csrf"),
} )
}
func NewDashboardHandler() DashboardHandler{
	return DashboardHandler{
		dashboardRepo: repository.NewDashboard(),
	}
}




