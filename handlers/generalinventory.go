package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)


type GeneralInventory struct {}

func NewGeneralInventory() GeneralInventory {
	return GeneralInventory{}
}
func(h * GeneralInventory)RenderGeneralInventoryPage(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/inventory/general/main", Data{
		"module": "General Inventory",
		"title":"General Inventory",
	})
}