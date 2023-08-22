package handlers

import (
	"fmt"
	"lift-fitness-gym/app/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

type MembershipPlanHandler struct {
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
		return c.JSON(http.StatusBadRequest,JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured",
		})
	}
	fmt.Println(plan)
	return nil
}



func NewMembershipPlanHandler() MembershipPlanHandler{
	return MembershipPlanHandler{}
}