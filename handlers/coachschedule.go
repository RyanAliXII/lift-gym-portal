package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)


type CoachSchedule struct {

}
func NewCoachScheduleHandler() CoachSchedule {
	return CoachSchedule{

	}
}
func( ctler * CoachSchedule)RenderCoachSchedulePage(c echo.Context) error {
	return c.Render(http.StatusOK, "coach/schedule/main", Data{
		"title": "Coach Schedule",
		"module": "Coach Schedule",
	})
}