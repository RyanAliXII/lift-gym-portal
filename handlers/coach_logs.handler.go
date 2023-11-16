package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/repository"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type CoachLogHandler struct {
	coachLogRepo repository.CoachLogRepository
	coachRepo repository.CoachRepository
}
func NewCoachLogHandler()CoachLogHandler{
	return CoachLogHandler{
		coachLogRepo: repository.NewCoachLogRepository(),
		coachRepo: repository.NewCoachRepository(),
	}
}
func (h *CoachLogHandler) RenderCoachLogPage(c echo.Context) error{

	contentType := c.Request().Header.Get("Content-Type")
	if contentType == "application/json"{
		logs, err := h.coachLogRepo.GetLogs()
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "GetLogsErr"))
		}
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: Data{
				"coachLogs": logs,
			},
			Message: "Coach logs fetched.",
		})
	}
	return c.Render(http.StatusOK, "admin/coach-logs/main", Data{
		 "title": "Coach Logs",
		 "module": "Coach Logs",
		 "csrf": c.Get("csrf"),
	} )
}
func (h * CoachLogHandler) NewLog(c echo.Context)error {
	log := model.CoachLog{}
	err := c.Bind(&log)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "bindErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured",
		})
	}
	err, fields := log.Validate()
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "validateErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": fields,
			},
			Message: "Validation error.",
		})
	}	
	err = h.coachLogRepo.NewLog(log)
	
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "getById"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Coach Logged",
	})
}


func (h * CoachLogHandler) UpdateLog(c echo.Context)error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "convErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured",
		})
	}
	log := model.CoachLog{}
	err = c.Bind(&log)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "bindErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured",
		})
	}
	log.Id = id
	err, fields := log.Validate()
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "validateErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": fields,
			},
			Message: "Validation error.",
		})
	}
	err = h.coachLogRepo.UpdateLog(log)

	if err != nil {
		logger.Error(err.Error(), zap.String("error", "getById"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured",
		})
	}
	

	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Coach Log Updated",
	})
}


func (h * CoachLogHandler)DeleteLog(c echo.Context)error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "convErr"))
		return c.JSON(http.StatusBadRequest,JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	err = h.coachLogRepo.DeleteLog(id)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "DeleteLogErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,		
			Message: "Unknown error occured",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Coach Logged",
	})
}


func (h * CoachLogHandler)Logout(c echo.Context)error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "convErr"))
		return c.JSON(http.StatusBadRequest,JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	err = h.coachLogRepo.LogoutCoach(id)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "LogoutCoach"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,		
			Message: "Unknown error occured",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Coach Logged",
	})
}



