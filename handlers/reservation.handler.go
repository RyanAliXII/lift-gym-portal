package handlers

import (
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)


type ReservationHandler struct {
	dateSlot repository.DateSlot
}
func NewReservationHandler () ReservationHandler{
	return ReservationHandler{
		dateSlot:  repository.NewDateSlotRepository(),
	}
}
func(h * ReservationHandler) RenderClientReservationPage(c echo.Context) error {
	return c.Render(http.StatusOK, "client/reservation/main", Data{})
}

func (h * ReservationHandler)GetDateSlots(c echo.Context) error {
	slots, err := h.dateSlot.GetSlots()
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "GetSlotsErr"))
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Data: Data{
			"slots" : slots,
		},
	})
}
