package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type MembershipPlanHandler struct {
	membershipPlanRepo repository.MembershipPlanRepository
}

func (h *MembershipPlanHandler) RenderMembershipPlanPage(c echo.Context) error{
	csrf := c.Get("csrf")
	return c.Render(http.StatusOK, "admin/membership-plan/main", Data{
		"csrf": csrf,
	} )
}
func (h *MembershipPlanHandler) NewMembershipPlan(c echo.Context) error{
	plan := model.MembershipPlan{}
	bindErr := c.Bind(&plan)
	validateErr, fieldErrs := plan.Validate()
	if validateErr != nil {
		return c.JSON(http.StatusBadRequest,JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured",
			Data: Data{
				"errors" : fieldErrs,
			},
		})
	}
	if bindErr != nil {
		logger.Error(bindErr.Error(), zap.String("error", "bindErr"))
		return c.JSON(http.StatusBadRequest,JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured",
		})
	}
	newPlanErr := h.membershipPlanRepo.New(plan)
	if newPlanErr != nil {
		logger.Error(newPlanErr.Error(), zap.String("error", "newPlanErr"))
		return c.JSON(http.StatusInternalServerError,JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Plan successfully created.",
	})
}



func NewMembershipPlanHandler() MembershipPlanHandler{
	return MembershipPlanHandler{
		membershipPlanRepo: repository.NewMembershipPlanRepository(),
	}
}