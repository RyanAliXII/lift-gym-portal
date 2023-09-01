package handlers

import (
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type MembershipRequestHandler struct {
	membershipPlanRepo repository.MembershipPlanRepository
}

func (h *MembershipRequestHandler) RenderClientMembershipRequest(c echo.Context) error{

	return c.Render(http.StatusOK, "client/membership-request/main", Data{
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


func NewMembershipRequestHandler() MembershipRequestHandler {
	return MembershipRequestHandler{
		 membershipPlanRepo: repository.NewMembershipPlanRepository(),
	}
}