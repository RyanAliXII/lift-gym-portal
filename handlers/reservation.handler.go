package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/app/repository"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)


type ReservationHandler struct {
	dateSlot repository.DateSlot
	timeSlot repository.TimeSlot
	reservation repository.Reservation
}
func NewReservationHandler () ReservationHandler{
	return ReservationHandler{
		dateSlot:  repository.NewDateSlotRepository(),
		timeSlot: repository.NewTimeSlotRepository(),
		reservation: repository.NewReservation(),
	}
}
func(h * ReservationHandler) RenderClientReservationPage(c echo.Context) error {
	contentType := c.Request().Header.Get("Content-Type")
	if contentType == "application/json" {
		sessionData := mysqlsession.SessionData{}
		sessionData.Bind(c.Get("sessionData"))
		reservations, err := h.reservation.GetClientReservation(sessionData.User.Id)
		if err != nil {
			logger.Error(err.Error(), zap.String("error","GetClientReservationErr"))
		}
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: Data{
				"reservations": reservations,
			},
			Message: "Fetch client's reservations.",
		})
	}
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
func (h * ReservationHandler)NewReservation(c echo.Context) error {
	reservation := model.Reservation{}
	err := c.Bind(&reservation)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "bindErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	fields, err := reservation.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors" : fields,
			},
			Message: "Validation error.",
		})
	}
	sessionData := mysqlsession.SessionData{}
	sessionData.Bind(c.Get("sessionData"))
	reservation.ClientId = sessionData.User.Id
	err = h.reservation.NewReservation(reservation)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "NewReservationErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Reservation created.",
	})
}
