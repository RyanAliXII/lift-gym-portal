package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type PackageRequestHandler struct {
}

func (r *PackageRequestHandler) RenderClientPackageRequestPage(c echo.Context) error {
	return c.Render(http.StatusOK, "client/package-request/main", Data{
		"csrf": c.Get("csrf"),
	})
}


func NewPackageRequestHandler() PackageRequestHandler {
	return PackageRequestHandler{}
}