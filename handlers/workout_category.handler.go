package handlers

import (
	"fmt"
	"lift-fitness-gym/app/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
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
func (h *WorkoutCategoryHandler) NewCategory(c echo.Context) error {
	category := model.WorkoutCateory{}
	err := c.Bind(&category)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "bindErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	err, fields := category.Validate()

	if err != nil {
		logger.Error(err.Error(), zap.String("error", "validationErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": fields,
			},
			Message: "Validation error.",
		})
	}
	fmt.Println(category)
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Category created.",
	})
}
func NewWorkoutCategoryHandler() WorkoutCategoryHandler {
	return WorkoutCategoryHandler{}
}