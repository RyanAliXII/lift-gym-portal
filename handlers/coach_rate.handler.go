package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type CoachRateHandler struct {
	coachRateRepo repository.CoachRateRepository
}

func NewCoachRateHandler() CoachRateHandler {

	return CoachRateHandler{
		coachRateRepo: repository.NewCoachRateRepository(),
	}
}

func (h *CoachRateHandler) RenderCoachRatePage(c echo.Context) error {

	
	return c.Render(http.StatusOK, "coach/rate/main", Data{
		"title":"Coach Rates",
		"module": "Coaching Rates",
		"csrf" : c.Get("csrf"),
	})
}

func (h *CoachRateHandler) NewRate(c echo.Context) error {
	rate := model.CoachRate{}
	err := c.Bind(&rate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	err, fields := rate.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": fields,
			},
			Message: "Validation error.",
		})
	}
	sessionData := mysqlsession.SessionData{}
	err = sessionData.Bind(c.Get("sessionData"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "bindSessionErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	rate.CoachId = sessionData.User.Id
	err = h.coachRateRepo.NewRate(rate)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "NewRateErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Coachin rate added.",
	})	
}