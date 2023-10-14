package handlers

import (
	"lift-fitness-gym/app/model"
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

func (h *CoachRateHandler) NewRate(c echo.Context) error {
	rate := model.CoachRate{}
	err := c.Bind(&rate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	err, fields := rate.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": fields,
			},
			Message: "Validation error.",
		})
	}
	
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Coachin rate added.",
	})	
}