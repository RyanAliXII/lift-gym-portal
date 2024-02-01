package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)


type GeneralInventory struct {
	generalInventoryRepo repository.GeneralInventory
}

func NewGeneralInventory() GeneralInventory {
	return GeneralInventory{
		generalInventoryRepo: repository.NewGeneralInventory(),
	}
}
func(h * GeneralInventory)RenderGeneralInventoryPage(c echo.Context) error {

	contentType := c.Request().Header.Get("Content-Type")
	if(contentType == "application/json"){
		items, err  := h.generalInventoryRepo.GetItems()
	
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "getItemsErr"))
		}
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: Data{
				"items": items,
			},
			Message: "New item added.",
		})
	}	
	return c.Render(http.StatusOK, "admin/inventory/general/main", Data{
		"module": "General Inventory",
		"title":"General Inventory",
	})
}
func (h  *  GeneralInventory)NewItem(c echo.Context) error {
	 body := model.GeneralItem{}
	 err := c.Bind(&body)
	 if err != nil {
		logger.Error(err.Error(), zap.String("error", "bindErr"))
	 }
	imageFile, err :=  c.FormFile("imageFile")
	if err != nil {
		if err.Error() != "http: no such file"{
			logger.Error(err.Error(), zap.String("error", "GetFileErr"))
		}
	}
	body.ImageFile = imageFile

	fieldsErrs, err := body.Validate()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data:Data{
				"errors": fieldsErrs,
			},
			Message: "Validation error.",
		})
	}
	
	err  = h.generalInventoryRepo.NewItem(body)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "NewItemErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Data: nil,
			Message: "Unknown error occured.",

		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Data: nil,
		Message: "New item added.",
	})
}