package handlers

import (
	"fmt"
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/app/pkg/objstore"
	"lift-fitness-gym/app/repository"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type WorkoutCategoryHandler struct {
	workoutCategoryRepo repository.WorkoutCategoryRepository
	clientRepo repository.ClientRepository
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
func (h *WorkoutCategoryHandler) RenderCoachCategoryPage(c echo.Context) error {

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
	return c.Render(http.StatusOK, "coach/workouts/category/main", Data{
		"csrf": c.Get("csrf"),
		"title": "Workout | Category",
		"module": "Workout Category",
	})	
}
func (h *WorkoutCategoryHandler) RenderClientWorkoutPage(c echo.Context)  error {

	sessionData := mysqlsession.SessionData{}
	bindErr := sessionData.Bind(c.Get("sessionData"))
	if bindErr != nil {
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
		   "message": "Unknown error occured",
	   })
	}
	client, getClientErr := h.clientRepo.GetById(sessionData.User.Id)
	
	if getClientErr != nil {
		logger.Error(getClientErr.Error(), zap.String("error", "getClientErr"))
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
		   "message": "Unknown error occured",
	   })
	}
	isInfoComplete := ((len(client.EmergencyContact) > 0) && (len(client.MobileNumber) > 0) && (len(client.Address) > 0))
	workoutCategories, _ := h.workoutCategoryRepo.GetCategories()
	return c.Render(http.StatusOK, "client/workouts/main", Data{
		"title": "Workouts",
		"module": "Workouts",
		"categories": workoutCategories,
		"client": client,
		"isInfoComplete": isInfoComplete,
	})

}
func (h *WorkoutCategoryHandler) RenderClientWorkoutsByCategoryId(c echo.Context)  error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "convErr"))
		c.Render(http.StatusNotFound,"partials/error/404-page", nil )
	}
	contentType := c.Request().Header.Get("Content-Type")

	
	sessionData := mysqlsession.SessionData{}
	bindErr := sessionData.Bind(c.Get("sessionData"))
	if bindErr != nil {
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
		   "message": "Unknown error occured",
	   })
	}
	client, getClientErr := h.clientRepo.GetById(sessionData.User.Id)
	if !client.IsMember || !client.IsVerified {
		return c.Render(http.StatusNotFound,"partials/error/404-page", nil )
	}
	if getClientErr != nil {
		logger.Error(getClientErr.Error(), zap.String("error", "getClientErr"))
		c.Render(http.StatusNotFound,"partials/error/404-page", nil )
	}
	workoutCategory, err := h.workoutCategoryRepo.GetCategoryById(id)
	fmt.Println(workoutCategory.Workouts)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "GetCategoryByIdError"))
		c.Render(http.StatusNotFound,"partials/error/404-page", nil )
	}

	if contentType == "application/json" {
		return c.JSON(http.StatusOK, JSONResponse{
			Data: Data{
			  "workouts": workoutCategory.Workouts,
			},
		})
	}
	title := fmt.Sprintf("Workout: %s", workoutCategory.Name)
	return c.Render(http.StatusOK, "client/workouts/id/main", Data{
		"title": title,
		"module": workoutCategory.Name,
		"category": workoutCategory,
		"publicURL" : objstore.PublicURL,
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

func (h *WorkoutCategoryHandler) UpdateCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "strConvErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	category := model.WorkoutCategory{}
	err = c.Bind(&category)
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
	category.Id  = id
	err = h.workoutCategoryRepo.UpdateCategory(category)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "newCategoryErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Category updated.",
	})
}

func (h *WorkoutCategoryHandler)DeleteCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "strConvErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	err = h.workoutCategoryRepo.DeleteCategory(id)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "newCategoryErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Category updated.",
	})
}
func NewWorkoutCategoryHandler() WorkoutCategoryHandler {
	return WorkoutCategoryHandler{
		workoutCategoryRepo: repository.NewWorkoutCategoryRepository(),
		clientRepo: repository.NewClientRepository(),
	}
}