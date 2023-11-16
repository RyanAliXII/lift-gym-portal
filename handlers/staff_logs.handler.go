package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/repository"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type StaffLogHandler struct {
	staffLogRepo repository.StaffLogRepository
	staffRepo repository.StaffRepository
}
func NewStaffLogHandler()StaffLogHandler{
	return StaffLogHandler{
		staffLogRepo: repository.NewStaffLogRepository(),
		staffRepo: repository.NewStaffRepository(),
	}
}
func (h *StaffLogHandler) RenderStaffLogPage(c echo.Context) error{

	contentType := c.Request().Header.Get("Content-Type")
	if contentType == "application/json"{
		logs, err := h.staffLogRepo.GetLogs()
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "GetLogsErr"))
		}
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: Data{
				"staffLogs": logs,
			},
			Message: "Staff logs fetched.",
		})
	}
	return c.Render(http.StatusOK, "admin/staff-logs/main", Data{
		 "title": "Staff Logs",
		 "module": "Staff Logs",
		 "csrf": c.Get("csrf"),
	} )
}
func (h * StaffLogHandler) NewLog(c echo.Context)error {
	log := model.StaffLog{}
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
	err = h.staffLogRepo.NewLog(log)
	
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "getById"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Staff Logged",
	})
}


func (h * StaffLogHandler) UpdateLog(c echo.Context)error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "convErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured",
		})
	}
	log := model.StaffLog{}
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
	err = h.staffLogRepo.UpdateLog(log)

	if err != nil {
		logger.Error(err.Error(), zap.String("error", "getById"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Staff Log Updated",
	})
}


func (h * StaffLogHandler)DeleteLog(c echo.Context)error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "convErr"))
		return c.JSON(http.StatusBadRequest,JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	err = h.staffLogRepo.DeleteLog(id)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "DeleteLogErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,		
			Message: "Unknown error occured",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Staff Logged",
	})
}


func (h * StaffLogHandler)Logout(c echo.Context)error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "convErr"))
		return c.JSON(http.StatusBadRequest,JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	err = h.staffLogRepo.LogoutStaff(id)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "LogoutStaff"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,		
			Message: "Unknown error occured",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Staff Logged",
	})
}



