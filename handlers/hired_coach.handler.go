package handlers

import (
	"fmt"
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/app/pkg/status"
	"lift-fitness-gym/app/repository"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type HiredCoachHandler  struct{
	hiredCoach repository.HiredCoachRepository
	client repository.ClientRepository
}

func NewHiredCoachHandler() HiredCoachHandler {
	return HiredCoachHandler{
		hiredCoach: repository.NewHiredCoachRepository(),
		client: repository.NewClientRepository(),
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
		Message: "Appointment Cancelled.",
	})
}

func (h *HiredCoachHandler) UpdateStatus(c echo.Context) error {
	id, err  := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error","atoiErrId"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	statusId, err := strconv.Atoi(c.QueryParam("statusId"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error","atoiErrStatus"))
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
	switch(statusId){
		case status.CoachAppointmentStatusApproved:
			 return h.handleApproval(c, id, statusId, sessionData.User.Id)
		case status.CoachAppointmentStatusPaid:
			 return h.handleMarkingAsPaid(c, id, statusId, sessionData.User.Id)
		case status.CoachAppointmentStatusCancelled:
			 return h.handleCancellation(c, id,  statusId, sessionData.User.Id)
		case status.CoachAppointmentStatusNoShow:
			return h.handleNoShow(c, id, statusId, sessionData.User.Id)
	}
	return c.JSON(http.StatusBadRequest, JSONResponse{
		Status: http.StatusBadRequest,
		Message: "Unknown action.",
	})
}

func (h *HiredCoachHandler)handleApproval(c echo.Context, id int, statusId int, coachId int) error {
	body := model.HiredCoach{}
	// err := c.Bind(&body)
	// if err != nil {
	// 	logger.Error(err.Error(), zap.String("error", "bindErr"))
	// 	return c.JSON(http.StatusBadRequest, JSONResponse{
	// 		Status: http.StatusBadRequest,
	// 		Message: "Unknown error occured.",
	// 	})
	// }

	// err, fields := body.ValidateMeetingTime()
	// if err != nil {
	// 	logger.Error(err.Error(), zap.String("error", "validation error"))
	// 	return c.JSON(http.StatusBadRequest, JSONResponse{
	// 		Status: http.StatusBadRequest,
	// 		Data: Data{
	// 			"errors": fields,
	// 		},
	// 	})
	// }
	body.Id = id
	body.StatusId = statusId
	body.CoachId = coachId
	err := h.hiredCoach.MarkAppointmentAsApproved(body)
	if err != nil{ 
		logger.Error(err.Error(), zap.String("error", "MarkAppointmentAsApproved"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Status updated.",
	})
}


func (h *HiredCoachHandler)handleMarkingAsPaid(c echo.Context, id int, statusId int, coachId int) error {
	body := model.HiredCoach{}
	body.Id = id
	body.StatusId = statusId
	body.CoachId = coachId
	err := h.hiredCoach.MarkAppointmentAsPaid(body)
	if err != nil{ 
		logger.Error(err.Error(), zap.String("error", "MarkAppointmentAsApproved"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Status updated.",
	})
}

func (h *HiredCoachHandler)handleCancellation(c echo.Context, id int, statusId int, coachId int) error {
	body := model.HiredCoach{}
	body.Remarks = c.FormValue("remarks")
	body.Id = id
	body.StatusId = statusId
	body.CoachId = coachId
	if(len(body.Remarks) > 0){
		body.Remarks = fmt.Sprintf("Cancelled by coach: %s ", body.Remarks )
	}else{
		body.Remarks = "Cancelled by coach"	
	}
	err := h.hiredCoach.CancelAppointment( body)
	if err != nil{ 
		logger.Error(err.Error(), zap.String("error", "MarkAppointmentAsCancelled"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Status updated.",
	})
}
func (h *HiredCoachHandler)handleNoShow(c echo.Context, id int, statusId int, coachId int) error {
	body := model.HiredCoach{}
	body.Id = id
	body.StatusId = statusId
	body.CoachId = coachId
err := h.hiredCoach.MarkAsNoShow(body)
	if err != nil{ 
		logger.Error(err.Error(), zap.String("error", "MarkAppointmentAsNoShow"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Status updated.",
	})
}


