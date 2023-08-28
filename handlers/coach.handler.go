package handlers

import (
	"lift-fitness-gym/app/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CoachHandler struct {
}

func (h *CoachHandler) RenderCoachPage(c echo.Context) error {
	csrf := c.Get("csrf")
	return c.Render(http.StatusOK, "admin/coach/main", Data{
		"title": "Coaches",
		"module": "Coaches",
		"csrf": csrf,
	})
}
func (h * CoachHandler)RenderCoachRegistrationPage(c echo.Context) error {
	csrf := c.Get("csrf")
	return c.Render(http.StatusOK, "admin/coach/register-coach", Data{
		"csrf": csrf,
		"title": "Coach | Registration",
		"module": "Registration Form",
	} )
}


func (h  *CoachHandler)NewCoach (c echo.Context) error {
	coach := model.Coach{}
	c.Bind(&coach)
	err, fieldErrs := coach.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": fieldErrs,
			},
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: nil,
			Message: "Coach has been registered.",
	})
}

func NewCoachHandler() CoachHandler{
	return CoachHandler{}
}