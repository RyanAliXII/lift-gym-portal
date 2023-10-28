package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)


type DateSlotHandler struct {}

func NewDateSlotHandler () DateSlotHandler{
	return DateSlotHandler{}
}
func(h * DateSlotHandler) RenderDateSlotPage (c echo.Context) error {
	return c.Render(http.StatusOK, "admin/reservation/date-slot/main", Data{
		"title": "Date Slots",
		"module": "Reservation Date Slots",
	})
}
