package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type MemberHandler struct {
}

func (h *MemberHandler) RenderMembersPage(c echo.Context) error{
	 csrf := c.Get("csrf")
	return c.Render(http.StatusOK, "admin/members/main", Data{
		"csrf": csrf,
	})
}
func NewMembersHandler() MemberHandler{
	return MemberHandler{}
}
