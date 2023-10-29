package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)


type TimeSlotHandler struct {
	timeSlotRepo repository.TimeSlot
}

func NewTimeSlotHandler () TimeSlotHandler {
	return TimeSlotHandler{
		timeSlotRepo:  repository.NewTimeSlotRepository(),
	}
}
func (h *TimeSlotHandler) RenderTimeSlotPage(c echo.Context) error {
	contentType := c.Request().Header.Get("Content-Type")
	if contentType == "application/json" {
		slots, err := h.timeSlotRepo.GetTimeSlots()
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "GetTimeSlotsErr"))
		}
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: Data{
				"slots": slots,
			},
		})
	}
	return c.Render(http.StatusOK, "admin/reservation/time-slot/main", Data{
		"title": "Time Slot",
		"module": "Reservation Time Slot",
	})
}
func (h *TimeSlotHandler) NewTimeSlot(c echo.Context) error {
	timeSlot := model.TimeSlot{}
	err := c.Bind(&timeSlot)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "bindErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	err, fields := timeSlot.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{Status: http.StatusBadRequest, Data: Data{
			"errors": fields,
		},
		Message: "Validation error.",
	})
	}
	err = h.timeSlotRepo.NewTimeSlot(timeSlot)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "newTimeSlotErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{Status: http.StatusInternalServerError, Message: "Unknown error occured."})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Time slot added.",
	})
}