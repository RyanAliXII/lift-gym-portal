package handlers

import (
	"fmt"
	"lift-fitness-gym/app/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type MemberHandler struct {


}


func (h *MemberHandler) RenderMembersPage(c echo.Context) error{
	 csrf := c.Get("csrf")
	return c.Render(http.StatusOK, "admin/members/main", Data{
		"csrf": csrf,
	})
}

func (h * MemberHandler)Subscribe(c echo.Context) error{
	subscribeBody := model.Subscribe{}
	bindErr := c.Bind(&subscribeBody)
	if bindErr != nil {
		logger.Error(bindErr.Error(), zap.String("error", "bindErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: nil,
			Message: "Unknown error  occured.",
		})
	}

	validatErr, _ := subscribeBody.Validate()
	if validatErr != nil {
		logger.Error(validatErr.Error(), zap.String("error", "validateErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: nil,
			Message: "Unknown error  occured.",
		})
	}
	fmt.Println(subscribeBody)
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Data: nil,
		Message: "Client subscribed.",
	})
}
func NewMembersHandler() MemberHandler{
	return MemberHandler{}
}
