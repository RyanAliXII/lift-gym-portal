package handlers

import (
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo/v4"
)


type CoachingPenalty struct {
	coachingPenaltyRepo repository.CoachingPenalty
}

func NewCoachingPenalty() CoachingPenalty {
	return CoachingPenalty{
		coachingPenaltyRepo: repository.NewCoachingPenalty(),
	}
}
func (h * CoachingPenalty) RenderPenaltyPage (c echo.Context) error{
	contentType := c.Request().Header.Get("content-type")
	if contentType == "application/json"{
		penalties, err := h.coachingPenaltyRepo.GetPenalties()
		if err != nil {
			logger.Error(err.Error())
		}
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: Data{
				"penalties": penalties,
			},
		})
	}
	return c.Render(http.StatusOK, "admin/coaching-penalty/main", Data{
			"title": "Coaching Penalty",
			"module": "Penalties",
	})
}