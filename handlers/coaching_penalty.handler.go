package handlers

import (
	"lift-fitness-gym/app/repository"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
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
func (h * CoachingPenalty) SettlePenalty (c echo.Context) error{
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "convErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	err = h.coachingPenaltyRepo.MarkAsSettled(id)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "markAsSettledErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Penalty settled.",
	})
}
func (h * CoachingPenalty)UnsettlePenalty(c echo.Context) error{
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "convErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	err = h.coachingPenaltyRepo.MarkAsUnSettled(id)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "markAsSettledErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Penalty settled.",
	})
}