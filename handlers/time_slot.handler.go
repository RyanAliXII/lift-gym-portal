package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)


type TimeSlotHandler struct {}

func NewTimeSlotHandler () TimeSlotHandler {
	return TimeSlotHandler{}
}

func (h *TimeSlotHandler) RenderTimeSlotPage(c echo.Context) error {

	return c.Render(http.StatusOK, "admin/reservation/time-slot/main", Data{
		"title": "Time Slot",
		"module": "Reservation Time Slot",
	})
}