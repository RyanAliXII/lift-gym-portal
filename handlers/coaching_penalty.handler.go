package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)


type CoachingPenalty struct {}

func NewCoachingPenalty() CoachingPenalty {
	return CoachingPenalty{}
}
func (h * CoachingPenalty) RenderPenaltyPage (c echo.Context) error{
	return c.Render(http.StatusOK, "admin/coaching-penalty/main", Data{
			"title": "Coaching Penalty",
			"module": "Penalties",
	})
}