package model

import "github.com/labstack/echo/v4"

type CoachAppointmentHandler struct{}

func NewCoachAppointmentHandler() CoachAppointmentHandler {
	return CoachAppointmentHandler{}
}

func (h *CoachAppointmentHandler) RenderClientAppointmentsPage(c echo.Context) error {


	return 
}