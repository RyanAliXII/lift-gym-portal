package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/repository"
	"net/http"

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
	
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Date slot/s has been added.",
	})
}
