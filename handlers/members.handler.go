package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type MemberHandler struct {
	memberRepository repository.MemberRepository 
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
	subscribeErr := h.memberRepository.Subscribe(subscribeBody)
	if subscribeErr != nil {
		logger.Error(subscribeErr.Error(), zap.String("error", "subscribeErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Data: nil,
			Message: "Unknown error occured.",

		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Data: nil,
		Message: "Client subscribed.",
	})
}
func NewMembersHandler() MemberHandler{
	return MemberHandler{
		memberRepository: repository.NewMemberRepository(),
	}
}
