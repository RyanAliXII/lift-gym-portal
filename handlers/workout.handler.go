package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type WorkoutHandler struct {
}

func (h *WorkoutHandler) RenderWorkoutPage(c echo.Context)  error {

	return c.Render(http.StatusOK, "admin/workouts/main", Data{
		"csrf" : c.Get("csrf"),
		"title": "Workouts",
		"module": "Workouts",
	})

}

func NewWorkoutHandler () WorkoutHandler {
	return WorkoutHandler{}
} 