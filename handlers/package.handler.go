package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/repository"
	"net/http"

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
	pkg := model.Package{}
	bindErr := c.Bind(&pkg)
	if bindErr != nil  {
		return c.JSON(http.StatusOK, Data{
			"status": http.StatusBadRequest,
			"message": "Unknown error occured.",

		})
	}
	newPackageErr := h.packageRepo.NewPackage(pkg)
	if newPackageErr != nil {
		return c.JSON(http.StatusOK, Data{
			"status": http.StatusInternalServerError,
			"message": "Unknown error occured.",

		})
	}
	return c.JSON(http.StatusOK, Data{
		"status": http.StatusOK,
		"message": "New package created.",

	})
}
func NewPackageHandler() PackageHandler {

	return PackageHandler{
		packageRepo: repository.NewPackageRepository(),
	}
}