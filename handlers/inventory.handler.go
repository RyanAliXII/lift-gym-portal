package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/repository"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type InventoryHandler struct {
	inventoryRepo repository.InventoryRepository
}
func NewInventoryHandler()InventoryHandler{
	return InventoryHandler{
		inventoryRepo: repository.NewInventoryRepository(),
	}
}
func (h *InventoryHandler) RenderInventoryPage(c echo.Context) error {
	contentType := c.Request().Header.Get("Content-Type")
	if contentType == "application/json" {
		equipments, err := h.inventoryRepo.GetEquipments()

		if err != nil{
			logger.Error(err.Error(),  zap.String("error","GetEquipmentsErr"))
		}
		stat, err := h.inventoryRepo.GetStat()
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "GetStatErr"))
		}
		return c.JSON(http.StatusOK, JSONResponse{Status: http.StatusOK, Data:Data{
			"equipments": equipments,
			"stat": stat,
		}})
	}
	return c.Render(http.StatusOK, "admin/inventory/main", Data{
		"title": "Equipment Inventory",
		"module": "Equipment Inventory",
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

	err = h.inventoryRepo.NewEquipment(equipment)
    if err != nil {
		logger.Error(err.Error(), zap.String("error", "NewEquipmentErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Data: nil,
		Message: "Equipment added.",
	})
}

func (h *InventoryHandler) UpdateEquipment(c echo.Context) error {
	equipment := model.Equipment{}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "strConvErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	err = c.Bind(&equipment)
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
	
	equipment.Id = id
	err = h.inventoryRepo.UpdateEquipment(equipment)
    if err != nil {
		logger.Error(err.Error(), zap.String("error", "NewEquipmentErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Data: nil,
		Message: "Equipment updated.",
	})
}

func (h *InventoryHandler) DeleteEquipment(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "strConvErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	err = h.inventoryRepo.DeleteEquipment(id)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "DeleteErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Data: nil,
		Message: "Equipment added.",
	})
}
