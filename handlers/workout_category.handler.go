package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type WorkoutCategoryHandler struct {
}

func (h *WorkoutCategoryHandler) RenderCategoryPage(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/workouts/category/main", Data{
		"csrf": c.Get("csrf"),
		"title": "Workout | Category",
		"module": "Workout Category",
	})	
}

func NewWorkoutCategoryHandler() WorkoutCategoryHandler {
	return WorkoutCategoryHandler{}
}