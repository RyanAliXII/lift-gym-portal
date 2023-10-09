package handlers

import (
	"context"
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/objstore"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type WorkoutHandler struct {
	objectStorage objstore.ObjectStorer
}

func (h *WorkoutHandler) RenderWorkoutPage(c echo.Context)  error {

	return c.Render(http.StatusOK, "admin/workouts/main", Data{
		"csrf" : c.Get("csrf"),
		"title": "Workouts",
		"module": "Workouts",
	})

}
func (h *WorkoutHandler) NewWorkout(c echo.Context)  error {
	workout := model.Workout{}
	workout.Name = c.FormValue("name")
	workout.Description = c.FormValue("description")
	// err := c.Bind(&workout)
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, JSONResponse{
	// 		Status: http.StatusBadRequest,
	// 		Message: "Unknown error occured",
	// 	})
	// }
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
		return c.JSON(http.StatusInternalServerError, JSONResponse{})
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
	}
} 