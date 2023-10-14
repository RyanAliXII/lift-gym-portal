package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ClientLogHandler struct {
	clientLogRepo repository.ClientLogRepository
	clientRepo repository.ClientRepository
}
func NewClientLogHandler()ClientLogHandler{
	return ClientLogHandler{
		clientLogRepo: repository.NewClientLogRepository(),
		clientRepo: repository.NewClientRepository(),
	}
}
func (h *ClientLogHandler) RenderClientLogPage(c echo.Context) error{

	contentType := c.Request().Header.Get("Content-Type")
	if contentType == "application/json"{
		logs, err := h.clientLogRepo.GetLogs()
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "GetLogsErr"))
		}
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: Data{
				"clientLogs": logs,
			},
			Message: "Client logs fetched.",
		})
	}
	return c.Render(http.StatusOK, "admin/client-logs/main", Data{
		 "title": "Client Logs",
		 "module": "Client Logs",
		 "csrf": c.Get("csrf"),
	} )
}
func (h * ClientLogHandler) NewLog(c echo.Context)error {
	log := model.ClientLog{}
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

	client, err := h.clientRepo.GetById(log.ClientId)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "getById"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured",
		})
	}
	
	log.IsMember = client.IsMember
	if log.IsMember {
		log.AmountPaid = 0
	}

	err = h.clientLogRepo.NewLog(log)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "NewLogerr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,		
			Message: "Unknown error occured",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Client Logged",
	})
}