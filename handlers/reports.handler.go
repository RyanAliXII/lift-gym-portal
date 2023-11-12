package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Report struct {}

func NewReportHandler() Report {
	return Report{}
}
func(h * Report) RenderReportPage(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/reports/main", Data{	
		"title":"Reports",
		"module": "Reports", 

	})
}