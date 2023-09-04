package handlers

import (
	"fmt"
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/app/pkg/status"
	"lift-fitness-gym/app/repository"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type PackageRequestHandler struct {
	packageRepo repository.PackageRepository
	packageRequestRepo repository.PackageRequestRepository
}

func (h *PackageRequestHandler) RenderClientPackageRequestPage(c echo.Context) error {

	contentType := c.Request().Header.Get("content-type")
	fmt.Println(contentType)
	if contentType == "application/json" {
		sessionData := mysqlsession.SessionData{}
		sessionData.Bind(c.Get("sessionData"))
		pkgRequests, err := h.packageRequestRepo.GetPackageRequestsByClientId(sessionData.User.Id)
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "getPackageRequestsByClientIdErr"))
			return c.JSON(http.StatusInternalServerError, JSONResponse{
				Status: http.StatusInternalServerError,
				Message: "Unknown error occured.",
			})
		}
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Message: "Package requests fetched.",
			Data: Data{
				"packageRequests": pkgRequests,
			},
		})
	}
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
			Status: http.StatusInternalServerError,
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
			Status: http.StatusBadRequest,
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
		Message: "Package request submitted.",
	})
}
func (h *PackageRequestHandler) UpdatePackageRequestStatus(c echo.Context)error {
	id,err := strconv.Atoi( c.Param("id"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "idConvertErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	statusId, err :=  strconv.Atoi(c.QueryParam("statusId"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "statusIdConvertErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	switch(statusId){
	case status.PackageRequestStatusCancelled:
		h.packageRequestRepo.CancelPackageRequest(id, statusId)
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Message: "Package status updated.",
		})
	default:
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown action.",
		})
	}
	
}
func NewPackageRequestHandler() PackageRequestHandler {
	return PackageRequestHandler{
		 packageRepo: repository.NewPackageRepository(),
		 packageRequestRepo: repository.NewPackageRequestRepository(),
	}
}