package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
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
	})
}

