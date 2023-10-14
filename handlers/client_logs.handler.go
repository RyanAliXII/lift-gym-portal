package handlers

import (
	"lift-fitness-gym/app/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ClientLogHandler struct {}
func NewClientLogHandler()ClientLogHandler{
	return ClientLogHandler{}
}
func (h *ClientLogHandler) RenderClientLogPage(c echo.Context) error{
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
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Client Logged",
	})
}