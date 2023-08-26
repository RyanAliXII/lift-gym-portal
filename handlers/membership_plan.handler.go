package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/repository"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type MembershipPlanHandler struct {
	membershipPlanRepo repository.MembershipPlanRepository
}

func (h *MembershipPlanHandler) RenderMembershipPlanPage(c echo.Context) error{
	csrf := c.Get("csrf")
	contentType := c.Request().Header.Get("content-type")
	if contentType == "application/json"{
		plans, getPlansErr := h.membershipPlanRepo.Get()
		if getPlansErr != nil {
			logger.Error(getPlansErr.Error(), zap.String("error", "getPlansErr"))
		}
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: Data{
				"membershipPlans": plans,
			},
		})
	}
	return c.Render(http.StatusOK, "admin/membership-plan/main", Data{
		"csrf": csrf,
		"title": "Memberships",
		"module": "Memberships",
	} )
}
func (h *MembershipPlanHandler) NewMembershipPlan(c echo.Context) error{
	plan := model.MembershipPlan{}
	bindErr := c.Bind(&plan)
	if bindErr != nil {
		logger.Error(bindErr.Error(), zap.String("error", "bindErr"))
		return c.JSON(http.StatusBadRequest,JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured",
		})
	}
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
func (h *MembershipPlanHandler) UpdatePlan(c echo.Context) error{
	plan := model.MembershipPlan{}
	id := c.Param("id")
	planId, convErr := strconv.Atoi(id)

	if convErr != nil {
		logger.Error(convErr.Error(), zap.String("error", "convErr"))
		return c.JSON(http.StatusBadRequest,JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured",
		})
	}
	bindErr := c.Bind(&plan)
	if bindErr != nil {
		logger.Error(bindErr.Error(), zap.String("error", "bindErr"))
		return c.JSON(http.StatusBadRequest,JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured",
		})
	}
	plan.Id = planId
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

	updatePlanErr := h.membershipPlanRepo.Update(plan)
	if updatePlanErr != nil {
		logger.Error(updatePlanErr.Error(), zap.String("error", "updatePlanErr"))
		return c.JSON(http.StatusInternalServerError,JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured",
		})

	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Membership plan updated.",
	})
}



func NewMembershipPlanHandler() MembershipPlanHandler{
	return MembershipPlanHandler{
		membershipPlanRepo: repository.NewMembershipPlanRepository(),
	}
}