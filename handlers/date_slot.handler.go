package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/repository"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)


type DateSlotHandler struct {
	dateSlotRepo repository.DateSlot
}

func NewDateSlotHandler () DateSlotHandler{
	return DateSlotHandler{
		dateSlotRepo: repository.NewDateSlotRepository(),
	}
}
func(h * DateSlotHandler) RenderDateSlotPage (c echo.Context) error {
	return c.Render(http.StatusOK, "admin/reservation/date-slot/main", Data{
		"title": "Date Slots",
		"module": "Reservation Date Slots",
	})
}

func(h * DateSlotHandler) NewSlot (c echo.Context) error {
	body := model.DateRangeBody{}
	err := c.Bind(&body)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "bindErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	err, fields := body.Validate()
	if err != nil {
		logger.Error(err.Error(), zap.String("error","validationErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": fields,
			},
			Message: "Validation error.",
		})
	}
    from, to, _ := body.ToTime()
	slots := h.toDateSlotModel(from, to)
	err = h.dateSlotRepo.NewSlots(slots)
	if err != nil {
        logger.Error(err.Error(), zap.String("error", "NewSlotsErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Date slot/s has been added.",
	})
}
func (h * DateSlotHandler) toDateSlotModel (from time.Time, to time.Time) []model.DateSlot{
	 numberOfDays := 1
	 OneDay := 24 //24 hrs
     duration := to.Sub(from)
	 if(duration.Hours() > 0){
	      numberOfDays = (int(duration.Hours()) / OneDay) + 1
	 }
	 slots := make([]model.DateSlot, 0)
	 for i := 0; i < numberOfDays; i++ {
		date := from.Add((time.Hour *  time.Duration(OneDay)) * time.Duration(i))
		slots = append(slots, model.DateSlot{
			Date: date.Format(time.DateOnly),
		})
	 }
	return slots
}