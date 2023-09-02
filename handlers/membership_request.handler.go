package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/app/pkg/status"
	"lift-fitness-gym/app/repository"
	"net/http"
	"strconv"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type MembershipRequestHandler struct {
	membershipPlanRepo repository.MembershipPlanRepository
	membershipRequestRepo repository.MembershipRequestRepository
}

func (h *MembershipRequestHandler) RenderClientMembershipRequest(c echo.Context) error{
	s, err := session.Get("client_sid", c)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, JSONResponse{
			Status: http.StatusUnauthorized,
			Data: nil,
			Message: "Unauthorized.",
		})
	}
	session := mysqlsession.SessionData{}
	session.Bind(s.Values["data"])
	contentType := c.Request().Header.Get("Content-Type")
	if contentType == "application/json" {
		requests, err := h.membershipRequestRepo.GetMembershipRequestsByClientId(session.User.Id)
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "getMembershipRequestsByClientIdErr"))
			return c.JSON(http.StatusInternalServerError, JSONResponse{
				Status: http.StatusInternalServerError,
				Message: "Unknown error occured",
			})
		}
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: Data{
				"membershipRequests" : requests,
			},
			Message: "Membership requests has been fetched.",
		})
	}
	return c.Render(http.StatusOK, "client/membership-request/main", Data{
		"csrf" :  c.Get("csrf"),
		"title": "Client | Membership Requests",
		"module": "Membership Requests",
	})
}


func (h *MembershipRequestHandler) RenderAdminMembershipRequest(c echo.Context) error{
	
	contentType := c.Request().Header.Get("Content-Type")
	if contentType == "application/json" {
		requests, err := h.membershipRequestRepo.GetMembershipRequests()
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "getMembershipRequestsByClientIdErr"))
			return c.JSON(http.StatusInternalServerError, JSONResponse{
				Status: http.StatusInternalServerError,
				Message: "Unknown error occured",
			})
		}
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: Data{
				"membershipRequests" : requests,
			},
			Message: "Membership requests has been fetched.",
		})
	}
	return c.Render(http.StatusOK, "admin/membership-request/main", Data{
		"csrf" :  c.Get("csrf"),
		"title": "Client | Membership Requests",
		"module": "Membership Requests",
	})
}


func (h * MembershipRequestHandler)GetUnrequestedMembershipPlans(c echo.Context) error{
	s, err := session.Get("client_sid", c)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, JSONResponse{
			Status: http.StatusUnauthorized,
			Data: nil,
			Message: "Unauthorized.",
		})
	}
	session := mysqlsession.SessionData{}
	session.Bind(s.Values["data"])
	plans, err := h.membershipPlanRepo.GetUnrequestedPlansOfClient(session.User.Id)

	if err != nil {
		logger.Error(err.Error(), zap.String("error", "getUnrequestedPlansOfClientErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Data: Data{
			"membershipPlans": plans,
		},
		Message: "Membership plans fetched.",
	})
}
func (h * MembershipRequestHandler) NewRequest(c echo.Context) error {
	request := model.MembershipRequest{}
	err := c.Bind(&request)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "bindErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occurred.",
			})
	}
	s, err := session.Get("client_sid", c)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, JSONResponse{
			Status: http.StatusUnauthorized,
			Data: nil,
			Message: "Unauthorized.",
		})
	}
	session := mysqlsession.SessionData{}
	session.Bind(s.Values["data"])
	request.ClientId = session.User.Id
	request.StatusId = status.MembershipRequestStatusPending

	err, _ = request.Validate()
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "validateErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occurred.",
			})
	}
	err = h.membershipRequestRepo.NewRequest(request)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "NewRequestErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occurred.",
		})

	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Request has been added.",
	})
}
func (h  *  MembershipRequestHandler) UpdateMembershipRequestStatus(c echo.Context) error {
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
	if statusId == status.MembershipRequestStatusCancelled {
		err := h.membershipRequestRepo.CancelMembershipRequest(id, "Cancelled by user.")
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "cancelMembershipRequestErr"))
			return c.JSON(http.StatusInternalServerError, JSONResponse{
				Status: http.StatusInternalServerError,
				Message: "Unknown error occured.",
			})
		}
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Message: "Membership request cancelled.",
		})
	}
	return c.JSON(http.StatusBadRequest, JSONResponse{
		Status: http.StatusBadRequest,
		Message: "Unknown action.",
	})
}

func NewMembershipRequestHandler() MembershipRequestHandler {
	return MembershipRequestHandler{
		 membershipPlanRepo: repository.NewMembershipPlanRepository(),
		membershipRequestRepo: repository.NewMembershipRequestRepository(),
	}
}