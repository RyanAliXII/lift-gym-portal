package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/objstore"
	"lift-fitness-gym/app/repository"
	"net/http"
	"strconv"

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
		"publicURL" : objstore.PublicURL,
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

func (h  *  GeneralInventory)UpdateItem(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "ConvertErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data:nil,
			Message: "Unknown error occured.",
		})
	}
	body := model.GeneralItem{}
	err = c.Bind(&body)
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
   body.Id = id
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
   
   err  = h.generalInventoryRepo.Update(body)
   if err != nil {
	   logger.Error(err.Error(), zap.String("error", "UpdateItemErr"))
	   return c.JSON(http.StatusInternalServerError, JSONResponse{
		   Status: http.StatusInternalServerError,
		   Data: nil,
		   Message: "Unknown error occured.",

	   })
   }
   return c.JSON(http.StatusOK, JSONResponse{
	   Status: http.StatusOK,
	   Data: nil,
	   Message: "Item updated.",
   })
}
func(h * GeneralInventory)DeleteItem(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "ConvertErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data:nil,
			Message: "Unknown error occured.",
		})
	}
	err = h.generalInventoryRepo.DeleteItem(id)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "DeleteItemErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Data: nil,
			Message: "Unknown error occured.",
 
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Data: nil,
		Message: "Item Deleted.",
	})
}