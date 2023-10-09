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
func (h *WorkoutHandler) NewWorkout(c echo.Context)  error {



	
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Workout has been added.",
	})

}

func NewWorkoutHandler () WorkoutHandler {
	return WorkoutHandler{}
} 