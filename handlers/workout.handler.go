package handlers

import (
	"context"
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/objstore"
	"lift-fitness-gym/app/repository"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type WorkoutHandler struct {
	objectStorage objstore.ObjectStorer
	workoutRepo repository.WorkoutRepository
}

func (h *WorkoutHandler) RenderWorkoutPage(c echo.Context)  error {
	contentType := c.Request().Header.Get("Content-Type")
	if contentType == "application/json"{
		workouts, err := h.workoutRepo.GetWorkouts()

		if err != nil {
			logger.Error(err.Error(), zap.String("error", "GetWorkoutsErr"))
		}
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: Data{
				"workouts": workouts,
			},
			Message: "Workouts fetched.",
		})
	}
	return c.Render(http.StatusOK, "admin/workouts/main", Data{
		"csrf" : c.Get("csrf"),
		"title": "Workouts",
		"module": "Workouts",
		"publicURL": objstore.PublicURL,  
	})

}
func (h *WorkoutHandler) NewWorkout(c echo.Context)  error {
	workout := model.Workout{}
	workout.Name = c.FormValue("name")
	workout.Description = c.FormValue("description")

	err, fields := workout.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": fields,
			},
			Message: "validation Error",
		})
	}
	file, err := c.FormFile("file")

	if err != nil {
		logger.Error(err.Error(), zap.String("error", "formFile"))
		fields["file"] = "Animated image is required."
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": fields,
			},
			Message: "validation Error",
		})
	}
	multiparFile, err := file.Open()
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "fileOpenErr"))
	    fields["file"] = "Animated image is required."
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": fields,
			},
			Message: "validation Error",
		})
	}
	id, err := uuid.NewUUID()
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "uuidNewErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	const folderName = "/workouts/images"
	fileKey, err := h.objectStorage.Upload(context.Background(), multiparFile, objstore.UploadConfig{
		FolderName: folderName,
		Filename: id.String(),
		AllowedFormats: []string{"jpg", "png", "gif"},
	})
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "uploadError"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	workout.ImagePath = fileKey
	err = h.workoutRepo.NewWorkout(workout)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "newWorkoutError"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Workout has been added.",
	})

}


func (h *WorkoutHandler) UpdateWorkout(c echo.Context)  error {
	workoutId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "strConvErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	
	workout := model.Workout{}
	workout.Name = c.FormValue("name")
	workout.Description = c.FormValue("description")
	
	err, fields := workout.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": fields,
			},
			Message: "validation Error",
		})
	}
	file, err := c.FormFile("file")

	if err != nil {
		logger.Error(err.Error(), zap.String("error", "formFile"))
		fields["file"] = "Animated image is required."
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": fields,
			},
			Message: "validation Error",
		})
	}
	multiparFile, err := file.Open()
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "fileOpenErr"))
	    fields["file"] = "Animated image is required."
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": fields,
			},
			Message: "validation Error",
		})
	}
	id, err := uuid.NewUUID()
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "uuidNewErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	oldWorkout, err := h.workoutRepo.GetWorkout(workoutId)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "getWorkoutErr"))
		return c.JSON(http.StatusNotFound, JSONResponse{
			Status: http.StatusNotFound,
			Message: "Not found",
		})
	}
	const folderName = "/workouts/images"
	fileKey, err := h.objectStorage.Upload(context.Background(), multiparFile, objstore.UploadConfig{
		FolderName: folderName,
		Filename: id.String(),
		AllowedFormats: []string{"jpg", "png", "gif"},
	})
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "uploadError"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	
	workout.ImagePath = fileKey
	workout.Id = workoutId
	err = h.workoutRepo.UpdateWorkout(workout)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "updateWorkoutErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	err = h.objectStorage.Remove(context.Background(), oldWorkout.ImagePath)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "removeErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Workout has been added.",
	})

}
func (h *WorkoutHandler) DeleteWorkout(c echo.Context)  error {
	workoutId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "strConvErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	oldWorkout, err := h.workoutRepo.GetWorkout(workoutId)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "getWorkoutErr"))
		return c.JSON(http.StatusNotFound, JSONResponse{
			Status: http.StatusNotFound,
			Message: "Not found",
		})
	}
	err = h.objectStorage.Remove(context.Background(), oldWorkout.ImagePath)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "removeErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	err = h.workoutRepo.DeleteWorkout(workoutId)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "deleteWorkoutErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Message: "Unknown error occured.",
			Status: http.StatusInternalServerError,
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Workout has been added.",
	})

}
func NewWorkoutHandler () WorkoutHandler {
	objstore, err := objstore.GetObjectStorage()
	if err != nil {
		panic(err)
	}
	return WorkoutHandler{
		objectStorage: objstore ,
		workoutRepo: repository.NewWorkoutRepository(),
	}
} 