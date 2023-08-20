package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/repository"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PackageHandler struct {
	packageRepo repository.PackageRepository
}

func (h *PackageHandler) RenderPackagePage(c echo.Context) error {
	pkgs := h.packageRepo.GetPackages()
	
	contentType := c.Request().Header.Get("Content-Type")
	if contentType == "application/json" {
		return c.JSON(http.StatusOK,  Data{
			"status": http.StatusOK,
			"data": Data{
				"packages": pkgs,
			},
			"message": "packages fetched successfully.",

		})
	}
	return c.Render(http.StatusOK, "packages/main", Data{
		"title": "Packages",
		"module":"Packages",
	})
}
func (h * PackageHandler) NewPackage (c echo.Context) error {
	description := c.FormValue("description")
	price := c.FormValue("price")
	parsedPrice, parseErr := strconv.ParseFloat(price, 64) 
	if parseErr != nil {
		return c.Redirect(http.StatusSeeOther, "/packages")
	}
	pkg := model.Package{
		Description: description,
		Price: parsedPrice,
	}
	newPackageErr := h.packageRepo.NewPackage(pkg)
	if newPackageErr != nil {
		return c.Redirect(http.StatusSeeOther, "/packages")
	}
	c.Set("sample", "hello world")
	return c.Redirect(http.StatusSeeOther, "/packages")
}
func NewPackageHandler() PackageHandler {

	return PackageHandler{
		packageRepo: repository.NewPackageRepository(),
	}
}