package handlers

import (
	"fmt"
	"lift-fitness-gym/app/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
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

func(h  * Report) CreateReport (c echo.Context) error {
	reportConfig := model.ReportConfig{}
	err := c.Bind(&reportConfig)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "bindErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	fmt.Println(reportConfig)
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Reports Generated.",
	})
}