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


func (h *HiredCoachHandler) RenderCoachAppointments(c echo.Context) error {
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
		appointments, err := h.hiredCoach.GetCoachAppointments(sessionData.User.Id)
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
				"appointments": appointments,
			},
		})
	}
	return c.Render(http.StatusOK,"coach/appointments/main", Data{
		"csrf" : c.Get("csrf"),
	})
}

func (h *HiredCoachHandler) CancelAppointmentByClient(c echo.Context) error {
	id, err  := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error","atoiErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	sessionData := mysqlsession.SessionData{}
	err = sessionData.Bind(c.Get("sessionData"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "sessionError"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status:http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	err = h.hiredCoach.CancelAppointmentByClient(model.HiredCoach{
		Id: id,
		ClientId: sessionData.User.Id,
		Remarks: "Cancelled by user.",
	})
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "cancel err"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Appoinment Cancelled.",
	})
}
