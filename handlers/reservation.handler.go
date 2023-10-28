package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)


type ReservationHandler struct {

}
func NewReservationHandler () ReservationHandler{
	return ReservationHandler{}
}
func(h * ReservationHandler) RenderClientReservationPage(c echo.Context) error {
	return c.Render(http.StatusOK, "client/reservation/main", Data{})
}
