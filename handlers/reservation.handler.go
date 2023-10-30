package handlers

import (
	"lift-fitness-gym/app/repository"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)


type ReservationHandler struct {
	dateSlot repository.DateSlot
	timeSlot repository.TimeSlot
}
func NewReservationHandler () ReservationHandler{
	return ReservationHandler{
		dateSlot:  repository.NewDateSlotRepository(),
		timeSlot: repository.NewTimeSlotRepository(),
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
func (h * ReservationHandler)GetTimeSlotsBasedOnDateSlotId(c echo.Context) error {
	dateSlotId, err := strconv.Atoi(c.Param("dateSlotId"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "ConvErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	slots, err := h.timeSlot.GetTimeSlotsBasedOnDateSlot(dateSlotId)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "GetTimeSlotBasedOnDateSlot"))
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Data: Data{
			"slots": slots,
		},
		Message: "Time slots fetched.",
	})
}

