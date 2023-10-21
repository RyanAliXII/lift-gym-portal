package handlers

import (
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type HiredCoachHandler  struct{
	hiredCoach repository.HiredCoachRepository
}

func NewHiredCoachHandler() HiredCoachHandler {
	return HiredCoachHandler{
		hiredCoach: repository.NewHiredCoachRepository(),
	}
}

func (h *HiredCoachHandler) RenderClientHiredCoachesPage(c echo.Context) error {
	contentType := c.Request().Header.Get("Content-Type")
	if contentType == "application/json" {
		
		sessionData := mysqlsession.SessionData{}
		err := sessionData.Bind(c.Get("sessionData"))
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "sessionError"))
			return c.JSON(http.StatusInternalServerError, JSONResponse{
				Status:http.StatusInternalServerError,
				Message: "Unknown error occured.",
			})
		}
		hiredCoahes, err := h.hiredCoach.GetCoachReservationByClientId(sessionData.User.Id)
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "GetCoachReservationByClientId"))
			return c.JSON(http.StatusInternalServerError, JSONResponse{
				Status:http.StatusInternalServerError,
				Message: "Unknown error occured.",
			})
		}
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: Data{
				"hiredCoaches": hiredCoahes,
			},
		})
	}
	return c.Render(http.StatusOK,"client/hired-coaches/main", Data{
		"csrf" : c.Get("csrf"),
	})
}