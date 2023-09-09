package handlers

import (
	"encoding/json"
	"fmt"
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/mailer"
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/app/pkg/objstore"
	"lift-fitness-gym/app/repository"
	"net/http"
	"path/filepath"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type ProfileHandler struct {
	clientRepo repository.ClientRepository
	verificationRepo  repository.VerificationRepository
	memberRepo repository.MemberRepository
	coachRepo repository.CoachRepository
	coachImage repository.CoachImageRepository
	objStorage objstore.ObjectStorer
}

func (h *ProfileHandler) RenderClientProfilePage(c echo.Context) error{
	csrf := c.Get("csrf")
	sessionData := mysqlsession.SessionData{}
	bindErr := sessionData.Bind(c.Get("sessionData"))
	if bindErr != nil {
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
		   "message": "Unknown error occured",
	   })
	}
	client, getClientErr := h.clientRepo.GetById(sessionData.User.Id)
	member , _ := h.memberRepo.GetMemberById(sessionData.User.Id)
	isMember := false
	if member.Id > 0 {
		isMember = true
	}
	var emailVerification model.EmailVerification

	if !client.IsVerified {
		emailVerification, _ = h.verificationRepo.GetLatestSentEmailVerification(client.Id)	
	}

	if getClientErr != nil {
		logger.Error(getClientErr.Error(), zap.String("error", "getClientErr"))
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
		   "message": "Unknown error occured",
	   })
	}
	c.Render(http.StatusOK, "client/profile/main", Data{
		"csrf": csrf,
		"title": "Profile",
		"module": "Profile",
		"profile": client,
		"emailVerification": emailVerification,
		"isMember": isMember,
		"member": member,
	})
	return nil
}
func (h * ProfileHandler) RenderCoachProfile (c echo.Context) error {
	sessionData := mysqlsession.SessionData{}
	sessionData.Bind(c.Get("sessionData"))
	coach, _ := h.coachRepo.GetCoachById(sessionData.User.Id)
	coachImages, err := h.coachImage.GetImagesByCoachId(sessionData.User.Id)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "get images error"))
	}
	imageURLS := []string{}
	for _, coachImage := range coachImages {
		url := fmt.Sprint(objstore.PublicURL,coachImage.Path)
		imageURLS = append(imageURLS, url)
	}
	imageBytes, _ := json.Marshal(imageURLS)
	return c.Render(http.StatusOK, "coach/profile/main", Data{
		"profile": coach,
		"csrf": c.Get("csrf"),
		"images": string(imageBytes),
	})
}
func (h * ProfileHandler)ChangeCoachPassword( c echo.Context) error {
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
func (h * ProfileHandler) CreateEmailVerification(c echo.Context) error {
	sessionData := mysqlsession.SessionData{}
    bindErr := sessionData.Bind(c.Get("sessionData"))
	if bindErr != nil {
		logger.Error(bindErr.Error(), zap.String("error",  "bindErr"))
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
		   "message": "Unknown error occured",
	   })
	}
	verification , createErr := h.verificationRepo.CreateEmailVerification(sessionData.User.Id)
	if createErr != nil {
		logger.Error(createErr.Error(), zap.String("error",  "createErr"))
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
		   "message": "Unknown error occured",
	   })
	}
	go mailer.SendEmailVerification([]string{sessionData.User.Email}, sessionData.User.GivenName, verification.PublicId)
	return c.JSON(http.StatusOK, Data{
		"status": http.StatusOK,
		"data": verification,
	   "message": "Email Verification Sent.",
   })
}
func (h  * ProfileHandler)ChangePassword (c echo.Context) error {
	sessionData := mysqlsession.SessionData{}
	sessionData.Bind(c.Get("sessionData"))
	client, err := h.clientRepo.GetByIdWithPassword(sessionData.User.Id)
	if err != nil {
		logger.Error(err.Error(), zap.String("error",  "GetByIdWithPasswordErr"))
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
		   "message": "Unknown error occured",
	   })
	}
	oldPassword := c.FormValue("oldPassword")
	err = validation.Validate(oldPassword, validation.Required, validation.Length(1, 0))
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				 "errors": Data{
					 "oldPassword": fmt.Sprint(err.Error(), "."),
				 },
			},
			Message: "Invalid old password value.",
		})
	}
	err = bcrypt.CompareHashAndPassword([]byte(client.Password), []byte(oldPassword))

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
	err = validation.Validate(newPassword, validation.Required, validation.Length(10, 30))

	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				 "errors": Data{
					 "newPassword": fmt.Sprint(err.Error(), "."),
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
	
	h.clientRepo.UpdatePassword(string(hashedNewPassword), client.Id)
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Data: nil,
		Message: "Password has been changed.",
	})
}
func(h * ProfileHandler)UpdatePublicProfile(c echo.Context) error {
	
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Message: "Unknown error occured",
			Status: http.StatusBadRequest,
		})
	}
	files := form.File["images"]
	sessionData := mysqlsession.SessionData{}
	sessionData.Bind(c.Get("sessionData"))
	images := make([]model.CoachImage, 0)
	for _, fileHeader := range files {
		multiPartFile, err := fileHeader.Open()
		defer multiPartFile.Close()
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "file open error."))
			return c.JSON(http.StatusInternalServerError, JSONResponse{
				Message: "Unknown error occured.",
				Status: http.StatusInternalServerError,
			})
		}
		id, err := uuid.NewUUID()
		
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "uuid err"))
			return c.JSON(http.StatusInternalServerError, JSONResponse{
				Status: http.StatusInternalServerError,
				Message: "Unknwon error occured",
			})

		}
	
		folderName := fmt.Sprintf("/coaches/images/%d/", sessionData.User.Id)
		fileExtension := filepath.Ext(fileHeader.Filename)
		fileName := fmt.Sprint(id.String(), fileExtension)
		err = h.objStorage.Upload(multiPartFile, folderName,  id.String())
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "upload error."))
			return c.JSON(http.StatusInternalServerError, JSONResponse{
				Message: "Unknown error occured.",
				Status: http.StatusInternalServerError,
			})
		}
		images = append(images, model.CoachImage{
			Path: fmt.Sprint(folderName,fileName),
			CoachId: sessionData.User.Id,
		})
			
	}
	err = h.coachImage.NewCoachImages(images)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "new coach image err"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Public profile updated.",
	})

}
func NewProfileHandler()ProfileHandler {
	objectStorage, _ := objstore.GetObjectStorage()
	return ProfileHandler{
		clientRepo: repository.NewClientRepository(),
		verificationRepo: repository.NewVerificationRepository() ,
		memberRepo: repository.NewMemberRepository(),
		coachRepo: repository.NewCoachRepository(),
		objStorage: objectStorage ,
		coachImage: repository.NewCoachImageRepository(),
	}
}