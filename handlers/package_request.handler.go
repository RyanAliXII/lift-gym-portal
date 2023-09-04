package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/app/pkg/status"
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type PackageRequestHandler struct {
	packageRepo repository.PackageRepository
	packageRequestRepo repository.PackageRequestRepository
}

func (h *PackageRequestHandler) RenderClientPackageRequestPage(c echo.Context) error {
	return c.Render(http.StatusOK, "client/package-request/main", Data{
		"csrf": c.Get("csrf"),
	})
}
func (h * PackageRequestHandler)GetUnrequestedPackages(c echo.Context) error{
	sessionData := mysqlsession.SessionData{}
	sessionData.Bind(c.Get("sessionData"))
	pkgs, err := h.packageRepo.GetUnrequestedPackageOfClient(sessionData.User.Id)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "getPkgsError"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusOK,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Data: Data{
			"packages": pkgs,
		},
	})

}

func(h * PackageRequestHandler)NewPackageRequest(c echo.Context) error {
	pkgRequest := model.PackageRequest{}
	err := c.Bind(&pkgRequest)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "bindErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusOK,
			Message: "Unknown error occured.",
		})
	}
	sessionData := mysqlsession.SessionData{}
	sessionData.Bind(c.Get("sessionData"))
	pkgRequest.ClientId = sessionData.User.Id
	pkgRequest.StatusId = status.PackageRequestStatusPending
	err = h.packageRequestRepo.NewPackageRequest(pkgRequest)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "NewPackageRequestErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Package has been successfully requested.",
	})
}
func NewPackageRequestHandler() PackageRequestHandler {
	return PackageRequestHandler{
		 packageRepo: repository.NewPackageRepository(),
		 packageRequestRepo: repository.NewPackageRequestRepository(),
	}
}