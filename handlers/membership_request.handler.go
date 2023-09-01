package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type MembershipRequestHandler struct {
}

func (h *MembershipRequestHandler) RenderClientMembershipRequest(c echo.Context) error{

	return c.Render(http.StatusOK, "client/membership-request/main", Data{
		"csrf" :  c.Get("csrf"),
		"title": "Client | Membership Requests",
		"module": "Membership Requests",
	})
}


func NewMembershipRequestHandler() MembershipRequestHandler {
	return MembershipRequestHandler{}
}