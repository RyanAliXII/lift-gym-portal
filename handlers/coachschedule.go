package handlers

import (
	"lift-fitness-gym/app/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)


type CoachSchedule struct {

}
func NewCoachScheduleHandler() CoachSchedule {
	return CoachSchedule{

	}
}
func( h * CoachSchedule)RenderCoachSchedulePage(c echo.Context) error {
	return c.Render(http.StatusOK, "coach/schedule/main", Data{
		"title": "Coach Schedule",
		"module": "Coach Schedule",
	})
}
func (h * CoachSchedule)NewSchedule(c echo.Context) error {

	schedule := model.CoachSchedule{}
	err := c.Bind(&schedule)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "bindErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: nil,
			Message: "Unknown error occured.",
		})
	}
	fieldErrs, err := schedule.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": fieldErrs,
			},
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Data: nil,
		Message: "New schedule created.",
	})
}