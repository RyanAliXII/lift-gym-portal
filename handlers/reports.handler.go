package handlers

import (
	"bytes"
	"fmt"
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/browser"
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/app/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/go-rod/rod/lib/proto"
	"github.com/labstack/echo/v4"
	"github.com/ysmood/gson"
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error())
		return c.Render(http.StatusNotFound, "partials/error/404-page", nil)
	}
	data, err := h.reportRepo.GetReportById(id)
	if err != nil{
		logger.Error(err.Error())
		return c.Render(http.StatusNotFound, "partials/error/404-page", nil)
	}
	data.WalkIns, err = h.reportRepo.GetWalkIns(data.StartDate, data.EndDate)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "getWalkInsErr"))
	}
	contentType := c.Request().Header.Get("content-type")
	if contentType == "application/json" {
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data:Data{
				"reportData": data,
			},
			Message: "Report data fetched.",
		})
	
	}
	return c.Render(http.StatusOK, "admin/reports/report-data", Data{
		"reportId": id,
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
	browser, err  := browser.NewBrowser()
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "newBrowserErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	sessionData := c.Get("sessionData")
	session := mysqlsession.SessionData{}
	session.Bind(sessionData)
	data, err := h.reportRepo.GenerateReportData(startDate, endDate, fmt.Sprintf("%s %s", session.User.GivenName, session.User.Surname))
	if err != nil {
		logger.Error(err.Error())
	}
	
	url := fmt.Sprintf("http://localhost/app/reports/%d", data.Id )
	page, err := browser.Goto(url)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "GotoErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	err = page.WaitStable(1 * time.Second)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "waitLoad"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	pdf, err := page.PDF(&proto.PagePrintToPDF{
		PaperWidth: gson.Num(8.5),
		PaperHeight: gson.Num(11),
	})
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "ToPDFErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	var buffer bytes.Buffer;
	_, err = buffer.ReadFrom(pdf)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "toBufferErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}

	return c.Stream(http.StatusOK, "application/pdf", &buffer)
}