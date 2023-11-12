package handlers

import (
	"fmt"
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Report struct {
	reportRepo repository.Report
}

func NewReportHandler() Report {
	return Report{
		reportRepo: repository.NewReport(),
	}
}
func(h * Report) RenderReportPage(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/reports/main", Data{	
		"title":"Reports",
		"module": "Reports", 

	})
}
func(h * Report) RenderReportData(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/reports/report-data", Data{	

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
	startDate,endDate, err := reportConfig.ToDateOnly()
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "toDateOnly"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	data, err := h.reportRepo.GenerateReportData(startDate, endDate)
	if err != nil {
		logger.Error(err.Error())
	}
	fmt.Println(data)
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Reports Generated.",
	})
}