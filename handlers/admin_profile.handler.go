package handlers

import (
	"context"
	"fmt"
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/app/pkg/objstore"
	"lift-fitness-gym/app/repository"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)



type AdminProfileHandler struct {
	userRepo repository.UserRepository
	objstore objstore.ObjectStorer
}

func NewAdminProfileHandler() AdminProfileHandler {
	store, _ := objstore.GetObjectStorage()
	return AdminProfileHandler{
		userRepo: repository.NewUserRepository(),
		objstore: store,
	}
}
func (h * AdminProfileHandler) RenderAdminProfile (c echo.Context) error {
	sessionData := mysqlsession.SessionData{}
	sessionData.Bind(c.Get("sessionData"))
	return c.Render(http.StatusBadRequest, "admin/profile/main",Data{})
}
func (h * AdminProfileHandler)ChangeAvatar (c echo.Context) error {
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

	result, err  := h.objstore.Upload(context.Background(), file, objstore.UploadConfig{
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
	err = h.userRepo.UpdateAvatar(sessionData.User.Id, result)
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

func (h * AdminProfileHandler) ChangePassword (c echo.Context) error {
	oldPassword := c.FormValue("oldPassword")
	newPassword := c.FormValue("newPassword")
	sessionData := mysqlsession.SessionData{}
	sessionData.Bind(c.Get("sessionData"))
	user, err := h.userRepo.GetUserById(sessionData.User.Id)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "GetUserByIdErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{Status: http.StatusBadRequest, 
			Data: Data{
			"errors":Data{
					"oldPassword" : "Old password is incorrect.",			
				},
			},
			Message: "Validation error."})	
	}


	err = validation.Validate(&newPassword, validation.Required.Error("New Password is required."), validation.Length(10, 30).Error("New password length must be atleast 10 characters to 30 characters."))
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{Status: http.StatusBadRequest, 
			Data: Data{
			"errors":Data{
					"newPassword" : err.Error(),			
				},
			},
			Message: "Validation error."})	
	}

	
	
	hashedPWD, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "hashingErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	err = h.userRepo.UpdateAdminPasswordByAccountId(user.AccountId, string(hashedPWD))
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "updatePWDErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "OK",
	})
}
func (h * AdminProfileHandler) GetAvatar (c echo.Context) error {
	
	sessionData := mysqlsession.SessionData{}
	sessionData.Bind(c.Get("sessionData"))
	avatarPath, err  := h.userRepo.GetUserAvatar(sessionData.User.Id)
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