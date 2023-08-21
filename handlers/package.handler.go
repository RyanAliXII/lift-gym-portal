package handlers

import (
	"fmt"
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
	csrf := c.Get("csrf")
	if contentType == "application/json" {
		return c.JSON(http.StatusOK,  Data{
			"status": http.StatusOK,
			"data": Data{
				"packages": pkgs,
			},
			"message": "packages fetched successfully.",

		})
	}
	return c.Render(http.StatusOK, "admin/packages/main", Data{
		"title": "Packages",
		"module":"Packages",
		"csrf": csrf,
	})
}
func (h * PackageHandler) NewPackage (c echo.Context) error {
	pkg := model.Package{}
	bindErr := c.Bind(&pkg)
	if bindErr != nil  {
		fmt.Println(bindErr.Error())
		return c.JSON(http.StatusBadRequest, Data{
			"status": http.StatusBadRequest,
			"message": "Unknown error occured.",

		})
	}
	newPackageErr := h.packageRepo.NewPackage(pkg)
	if newPackageErr != nil {
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
			"message": "Unknown error occured.",

		})
	}
	return c.JSON(http.StatusOK, Data{
		"status": http.StatusOK,
		"message": "Package created.",

	})
}

func (h  PackageHandler) UpdatePackage(c echo.Context) error {

	
	id, convErr :=  strconv.Atoi(c.Param("id"))
	pkg := model.Package{}
	bindErr := c.Bind(&pkg)
	if bindErr != nil  {
		fmt.Println(bindErr.Error())
		return c.JSON(http.StatusBadRequest, Data{
			"status": http.StatusBadRequest,
			"message": "Unknown error occured.",

		})
	}
	if convErr != nil || pkg.Id != id {
		fmt.Println("Conversion error")
		return c.JSON(http.StatusBadRequest, Data{
			"status": http.StatusBadRequest,
			"message": "Unknown error occured.",

		})
	}
	updatePackageErr := h.packageRepo.UpdatePackage(pkg)
	if updatePackageErr != nil {
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
			"message": "Unknown error occured.",

		})
	}
	return c.JSON(http.StatusOK, Data{
		"status": http.StatusOK,
		"message": "Package updated.",

	})

}
func NewPackageHandler() PackageHandler {
	return PackageHandler{
		packageRepo: repository.NewPackageRepository(),
	}
}