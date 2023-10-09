package handlers

import (
	"lift-fitness-gym/app/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type InventoryHandler struct {
}


func NewInventoryHandler()InventoryHandler{
	return InventoryHandler{}
}
func (h *InventoryHandler) RenderInventoryPage(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/inventory/main", Data{
		"title": "Equipment Inventory",
		"module": "Inventory",
		"csrf": c.Get("csrf"),
	})
}
func (h *InventoryHandler) NewEquipment(c echo.Context) error {
	equipment := model.Equipment{}
	err := c.Bind(&equipment)
	
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "bindErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}

	err, fields := equipment.Validate() 
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "validationErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": fields,
			},
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Data: nil,
		Message: "Equipment added.",
	})
}
