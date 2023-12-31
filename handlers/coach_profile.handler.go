package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/app/pkg/objstore"
	"lift-fitness-gym/app/repository"
	"net/http"
	"slices"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type CoachProfileHandler struct{
	coachImage repository.CoachImageRepository
	objStorage objstore.ObjectStorer
	coachRepo repository.CoachRepository
}
func (h * CoachProfileHandler) RenderProfilePage (c echo.Context) error {
	sessionData := mysqlsession.SessionData{}
	sessionData.Bind(c.Get("sessionData"))
	coach, _ := h.coachRepo.GetCoachById(sessionData.User.Id)
	coachImages, err := h.coachImage.GetImagesByCoachId(sessionData.User.Id)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "get images error"))
	}
	imageURLS := []string{}
	for _, coachImage := range coachImages {
		url := fmt.Sprint(objstore.PublicURL,"/",coachImage.Path)
		imageURLS = append(imageURLS, url)
	}
	imageBytes, _ := json.Marshal(imageURLS)
	return c.Render(http.StatusOK, "coach/profile/main", Data{
		"profile": coach,
		"csrf": c.Get("csrf"),
		"images": string(imageBytes),
	})
}
func (h * CoachProfileHandler)ChangePassword( c echo.Context) error {
	sessionData := mysqlsession.SessionData{}
	err := sessionData.Bind(c.Get("sessionData"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, JSONResponse{
			Status: http.StatusUnauthorized,
			Message: "Unauthorized.",
		})
	}
	coach, _ := h.coachRepo.GetCoachByIdWithPassword(sessionData.User.Id)
	oldPassword := c.FormValue("oldPassword")
	err = validation.Validate(oldPassword, validation.Required.Error("Old password is required."), validation.Length(1, 0).Error("Old password is required."))
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				 "errors": Data{
					 "oldPassword": err.Error(),
				 },
			},
			Message: "Invalid old password value.",
		})
	}
	err = bcrypt.CompareHashAndPassword([]byte(coach.Password), []byte(oldPassword))
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
		   Status: http.StatusBadRequest,
		   Data: Data{
				"errors": Data{
					"oldPassword": "Old password is incorrect.",
				},
		   },
		   Message: "Old password is incorrect.",
	   })

	}

	newPassword := c.FormValue("newPassword")
	err = validation.Validate(newPassword, validation.Required.Error("New password is required."), validation.Length(10, 30).Error("New password must be 10 to 30 characters."))
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				 "errors": Data{
					 "newPassword": err.Error(),
				 },
			},
			Message: "Invalid new password value.",
		})
	}
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(err.Error(), zap.String("error",  "generatePassword"))
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
		   "message": "Unknown error occured",
	   })
	}
	err = h.coachRepo.UpdatePassword(string(hashedNewPassword), coach.Id)
	if err != nil {
		logger.Error(err.Error(), zap.String("error",  "UpdatePasswordErr"))
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
		   "message": "Unknown error occured",
	   })
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Password has been changed.",
	})
}
func (h * CoachProfileHandler) UpdatePublicProfile(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Message: "Unknown error occured",
			Status:  http.StatusBadRequest,
		})
	}
	files := form.File["images"]
	sessionData := mysqlsession.SessionData{}
	sessionData.Bind(c.Get("sessionData"))

	//get the list of coach images from database 
	storedImagesPath, err := h.coachImage.GetImagesPathByCoachId(sessionData.User.Id)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "GetImagesPathByCoachId"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Message: "Unknown error occured",
			Status:  http.StatusInternalServerError,
		})
	}

	folderName := fmt.Sprintf("coaches/images/%d/", sessionData.User.Id)
	uploadedImagesMap := map[string]string{}

	const NumberOfFilesAllowed = 3
	const FiveMegabytes = 5000000

	if len(files) > NumberOfFilesAllowed {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Message: "Files uploaded exceeds the number of files allowed.",
			Status:  http.StatusBadRequest,
		})
	}
	
	// loop through files of form data
	uploadedFilesNewNameMap := map[string]string{}
	for _, fileHeader := range files {
		if fileHeader.Size > FiveMegabytes {
			logger.Error("File exceeds the allow sized of 5MB. File will be skipped.", zap.Int64("fileSize", fileHeader.Size),)
			continue
		}
		multiPartFile, err := fileHeader.Open()
		defer multiPartFile.Close()
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "file open error."))
			return c.JSON(http.StatusInternalServerError, JSONResponse{
				Message: "Unknown error occured.",
				Status:  http.StatusInternalServerError,
			})
		}

		id, err := uuid.NewUUID()
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "uuid err"))
			return c.JSON(http.StatusInternalServerError, JSONResponse{
				Status:  http.StatusInternalServerError,
				Message: "Unknown error occured",
			})

		}
		//check if images is already been uploaded, if uploaded, just skip it.
		uploadedFilename := fileHeader.Filename
		uploadedFileFullPath := fmt.Sprint(folderName, uploadedFilename)
		uploadedImagesMap[uploadedFileFullPath] = fileHeader.Filename
		newFilename := id.String()
		if slices.Contains(storedImagesPath, uploadedFileFullPath) {
			continue
		}
		uploadedFilesNewNameMap[uploadedFilename] = newFilename
		ctx, cancel := context.WithCancel(context.Background())
		var fileIdChan = make(chan string)
		//sender function for uploading
		go func(channel chan <- string){
			publicId, err := h.objStorage.Upload(ctx,multiPartFile, objstore.UploadConfig{
				FolderName: folderName,
				Filename: newFilename,
				AllowedFormats: []string{"jpg", "png", "webp", "jpeg"},
			})
			if err != nil {
				logger.Error(err.Error(), zap.String("error", "UploadErr"))
				cancel()
				close(fileIdChan)
				return
			}
			channel <- publicId
			close(fileIdChan)
		}(fileIdChan)
		
		//receiver function for uploading
		go func(channel <-chan string) {
			select {
				case <-ctx.Done():
						return
				case fileId := <- channel:
					err := h.coachImage.NewCoachImage(model.CoachImage{
						Path: fileId,
						CoachId: sessionData.User.Id,
					})
					if err != nil {
							logger.Error(err.Error(), zap.String("error","NewCoachImageErr"))
					}
				}
		}(fileIdChan)
	}
	for _, storedImagePath := range storedImagesPath {
		filename := uploadedImagesMap[storedImagePath]
		if filename != "" {
			continue
		}
		ctx, cancel := context.WithCancel(context.Background())
		var fileIdChan = make(chan string)

		//sender function for deletion
		go func(channel chan <- string, imagePath string){
			err := h.objStorage.Remove(ctx, imagePath)
			if err != nil {
				logger.Error(err.Error(), zap.String("error", "Remove"))
				cancel()
				close(fileIdChan)
				return
			}
			fileIdChan <- imagePath
			close(fileIdChan)
		}(fileIdChan, storedImagePath)

		//receiver function for deletion
		go func (channel <- chan string){
			select {
				case <-ctx.Done():
					return
				case fileId := <- channel:
					err := h.coachImage.DeleteCoachImage(model.CoachImage{
						Path: fileId,
						CoachId: sessionData.User.Id,
				})
				if err != nil {
					logger.Error(err.Error(), zap.String("error", "Remove"))
					cancel()
				}
			}
		}(fileIdChan)

	}
	err = h.coachRepo.UpdateCoachDescription(sessionData.User.Id, c.FormValue("description"))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "UpdateCoachDescriptionErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status:  http.StatusOK,
		Data: Data{
			"uploadedFilesNewName": uploadedFilesNewNameMap,
		},
		Message: "Public profile updated.",
	})

}

func (h * CoachProfileHandler)ChangeAvatar (c echo.Context) error {
	fileHeader, err := c.FormFile("filepond")
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	file, err := fileHeader.Open()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	defer file.Close()
	sessionData := mysqlsession.SessionData{}
	sessionData.Bind(c.Get("sessionData"))
	uuid := uuid.New()

	result, err  := h.objStorage.Upload(context.Background(), file, objstore.UploadConfig{
		FolderName: "avatars",
		Filename: uuid.String(),
		AllowedFormats: []string{"jpg", "png", "webp", "jpeg"},
	} )

	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	err = h.coachRepo.UpdateAvatar(sessionData.User.Id, result)
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Avatar changed successfully",
	})
}
func (h * CoachProfileHandler) GetAvatar (c echo.Context) error {
	sessionData := mysqlsession.SessionData{}
	sessionData.Bind(c.Get("sessionData"))
	avatarPath, err  := h.coachRepo.GetUserAvatar(sessionData.User.Id)
	if err != nil{
		logger.Error(err.Error())
	}
	avatarUrl := ``
	if(len(avatarPath) == 0){
		avatarUrl = fmt.Sprintf("https://ui-avatars.com/api/?name=%s+%s", sessionData.User.GivenName, sessionData.User.Surname)
	}else{
	  avatarUrl  = fmt.Sprintf("%s/%s", objstore.PublicURL, avatarPath)
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Data: Data{
			"avatarUrl": avatarUrl,
		},
		Message: "Avatar fetched.",
	})
} 


func NewCoachProfileHandler() CoachProfileHandler {
	objectStorage, _ := objstore.GetObjectStorage()
	return CoachProfileHandler{
		coachImage: repository.NewCoachImageRepository(),
		objStorage: objectStorage,
		coachRepo: repository.NewCoachRepository(),
	}
}