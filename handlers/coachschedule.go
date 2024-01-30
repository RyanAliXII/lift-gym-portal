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


type CoachSchedule struct {
	coachSchedRepo repository.CoachSchedule
}
func NewCoachScheduleHandler() CoachSchedule {
	return CoachSchedule{
		coachSchedRepo: repository.NewCoachSchedule(),
	}
}
func( h * CoachSchedule)RenderCoachSchedulePage(c echo.Context) error {
	contentType := c.Request().Header.Get("Content-Type")
	if(contentType == "application/json"){
		session := mysqlsession.SessionData{}
		session.Bind(c.Get("sessionData"))
		scheds, err := h.coachSchedRepo.GetSchedulesByCoachId(session.User.Id)
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "GetSchedulesByCoachIdErr"))
		}
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: Data{
				"schedules": scheds,
			},
			Message: "",
		})
	}
	return c.Render(http.StatusOK, "coach/schedule/main", Data{
		"title": "Coach Schedule",
		"module": "Appointment Schedules",
	})
}
func (h * CoachSchedule)NewSchedule(c echo.Context) error {
	schedule := model.CoachSchedule{}
	err := c.Bind(&schedule)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "bindErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: nil,
			Message: "Unknown error occured.",
		})
	}
	fieldErrs, err := schedule.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": fieldErrs,
			},
			Message: "Unknown error occured.",
		})
	}
	session := mysqlsession.SessionData{}
	session.Bind(c.Get("sessionData"))
	schedule.CoachId = session.User.Id
	err = h.coachSchedRepo.NewSchedule(schedule)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "NewScheduleErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Data: nil,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Data: nil,
		Message: "New schedule created.",
	})
}
func (h * CoachSchedule)UpdateSchedule(c echo.Context) error {
	id, err  := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "convertErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: nil,
			Message: "Unknown error occured.",
		})
	}
	schedule := model.CoachSchedule{}
	err = c.Bind(&schedule)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "bindErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: nil,
			Message: "Unknown error occured.",
		})
	}
	fieldErrs, err := schedule.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": fieldErrs,
			},
			Message: "Unknown error occured.",
		})
	}
	session := mysqlsession.SessionData{}
	session.Bind(c.Get("sessionData"))
	schedule.CoachId = session.User.Id
	schedule.Id = id
	err = h.coachSchedRepo.UpdateSchedule(schedule)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "NewScheduleErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Data: nil,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Data: nil,
		Message: "Schedule updated.",
	})
}

func (h * CoachSchedule)DeleteSchedule(c echo.Context) error {
	id, err  := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "convertErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: nil,
			Message: "Unknown error occured.",
		})
	}
	schedule := model.CoachSchedule{}
	session := mysqlsession.SessionData{}
	session.Bind(c.Get("sessionData"))
	schedule.CoachId = session.User.Id
	schedule.Id = id
	err = h.coachSchedRepo.DeleteSchedule(schedule)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "NewScheduleErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Data: nil,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Data: nil,
		Message: "Schedule deleted.",
	})
}