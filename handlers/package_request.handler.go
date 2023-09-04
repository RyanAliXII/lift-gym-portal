package handlers

import (
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type PackageRequestHandler struct {
	packageRepo repository.PackageRepository
}

func (r *PackageRequestHandler) RenderClientPackageRequestPage(c echo.Context) error {
	return c.Render(http.StatusOK, "client/package-request/main", Data{
		"csrf": c.Get("csrf"),
	})
}
func (r * PackageRequestHandler)GetUnrequestedPackages(c echo.Context) error{
	sessionData := mysqlsession.SessionData{}
	sessionData.Bind(c.Get("sessionData"))
	pkgs, err := r.packageRepo.GetUnrequestedPackageOfClient(sessionData.User.Id)

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
func NewPackageRequestHandler() PackageRequestHandler {
	return PackageRequestHandler{
		 packageRepo: repository.NewPackageRepository(),
	}
}