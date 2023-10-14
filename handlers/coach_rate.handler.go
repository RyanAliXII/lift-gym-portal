package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type CoachRateHandler struct {
}

func NewCoachRateHandler() CoachRateHandler {

	return CoachRateHandler{}
}

func (h *CoachRateHandler) RenderCoachRatePage(c echo.Context) error {

	
	return c.Render(http.StatusOK, "coach/rate/main", Data{
		"title":"Coach Rates",
		"module": "Coaching Rates",
		"csrf" : c.Get("csrf"),
	})
}