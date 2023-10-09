package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type WorkoutCategoryHandler struct {
	workoutCategoryRepo repository.WorkoutCategoryRepository
}

func (h *WorkoutCategoryHandler) RenderCategoryPage(c echo.Context) error {

	contentType := c.Request().Header.Get("Content-Type")
	if contentType == "application/json"  {
		categories, err := h.workoutCategoryRepo.GetCategories()
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "GetCategoriesErr"))
		}
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: Data{
				"categories": categories,
			},
		})
	}
	return c.Render(http.StatusOK, "admin/workouts/category/main", Data{
		"csrf": c.Get("csrf"),
		"title": "Workout | Category",
		"module": "Workout Category",
	})	
}
func (h *WorkoutCategoryHandler) NewCategory(c echo.Context) error {
	category := model.WorkoutCategory{}
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
	err = h.workoutCategoryRepo.NewCategory(category)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "newCategoryErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Category created.",
	})
}
func NewWorkoutCategoryHandler() WorkoutCategoryHandler {
	return WorkoutCategoryHandler{
		workoutCategoryRepo: repository.NewWorkoutCategoryRepository(),
	}
}