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

type CoachRateHandler struct {
	coachRateRepo repository.CoachRateRepository
}

func NewCoachRateHandler() CoachRateHandler {

	return CoachRateHandler{
		coachRateRepo: repository.NewCoachRateRepository(),
	}
}

func (h *CoachRateHandler) RenderCoachRatePage(c echo.Context) error {

	contentType := c.Request().Header.Get("Content-Type")
	if contentType == "application/json"{
		sessionData := mysqlsession.SessionData{}
		sessionData.Bind(c.Get("sessionData"))
		rates, err := h.coachRateRepo.GetRatesByCoachId(sessionData.User.Id)
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "GetRatesByCoachId"))
		}
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: Data{
				"rates": rates,
			},
			Message: "Coach rates fetched.",
		})
	}
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

func (h *CoachRateHandler)UpdateRate(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	rate := model.CoachRate{}
	err = c.Bind(&rate)
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
	rate.Id = id
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
	err = h.coachRateRepo.UpdateRate(rate)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "UpdateRateErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Coachin rate updated.",
	})	
}