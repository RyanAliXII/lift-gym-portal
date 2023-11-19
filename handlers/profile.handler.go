package handlers

import (
	"context"
	"fmt"
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/mailer"
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/app/pkg/objstore"
	"lift-fitness-gym/app/repository"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nyaruka/phonenumbers"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type ProfileHandler struct {
	clientRepo repository.ClientRepository
	verificationRepo  repository.VerificationRepository
	memberRepo repository.MemberRepository
	coachRepo repository.CoachRepository
	objstore objstore.ObjectStorer
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
	member , err := h.memberRepo.GetMemberById(sessionData.User.Id)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "GetMemberByIdErr"))
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
		"isMember": client.IsMember,
		"member": member,
	})
	return nil
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

func(h * ProfileHandler) UpdateProfile (c echo.Context) error {
	mobileNumber := c.FormValue("mobileNumber")
	emergencyContact := c.FormValue("emergencyContact")
	address := c.FormValue("address")
	session := mysqlsession.SessionData{}
	session.Bind(c.Get("sessionData"))
	if len(mobileNumber) > 0 {
		return h.handleMobileNumber(c, mobileNumber, session)
	}
	if len(emergencyContact) > 0 {
		return h.handleEmergencyContact(c, emergencyContact, session)
	}
	if len(address) > 0 {
		return h.handleAddress(c, address, session)
	}
	return c.JSON(http.StatusOK, JSONResponse{Status: http.StatusOK, Message: "OK"})
}


func(h * ProfileHandler) handleMobileNumber(c echo.Context, mobileNumber string, session mysqlsession.SessionData) error{
	err := validation.Validate(mobileNumber, validation.By(ValidateMobile))
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": Data{
					"mobileNumber": err.Error(),
				},
			},
		})
	}
	err = h.clientRepo.UpdateMobileNumberOnce(session.User.Id, mobileNumber)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "UpdateMobileNumber"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{Status: http.StatusOK, Message: "Mobile number updated."})
}
func(h * ProfileHandler) handleEmergencyContact(c echo.Context, emergencyContact string, session mysqlsession.SessionData) error{
	err := validation.Validate(emergencyContact, validation.By(ValidateMobile))
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": Data{
					"emergencyContact": err.Error(),
				},
			},
		})
	}
	err = h.clientRepo.UpdateEmergencyContactOnce(session.User.Id, emergencyContact)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "UpdateEmergencyContact"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{Status: http.StatusOK, Message: "Emergency contact updated."})
}


func(h * ProfileHandler) handleAddress(c echo.Context, address string, session mysqlsession.SessionData) error{
	err := h.clientRepo.UpdateAddressOnce(session.User.Id, address)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "UpdateAddress"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{Status: http.StatusOK, Message: "Emergency contact updated."})
}

func ValidateMobile (value interface {}) error {
	mobileNumber, _ := value.(string)
	fmt.Println(mobileNumber)
	p, _ := phonenumbers.Parse(mobileNumber, "PH")
	isValid := phonenumbers.IsValidNumberForRegion(p, "PH")
	if !isValid {
		return fmt.Errorf("Invalid number")
	}
	return nil
}

func (h * ProfileHandler)ChangeAvatar(c echo.Context) error {
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
	err = h.clientRepo.UpdateAvatar(sessionData.User.Id, result)
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

func (h * ProfileHandler) GetAvatar (c echo.Context) error {
	sessionData := mysqlsession.SessionData{}
	sessionData.Bind(c.Get("sessionData"))
	avatarPath, err  := h.clientRepo.GetUserAvatar(sessionData.User.Id)
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

func NewProfileHandler()ProfileHandler {
	objstore, _ := objstore.GetObjectStorage()
	return ProfileHandler{
		clientRepo: repository.NewClientRepository(),
		verificationRepo: repository.NewVerificationRepository() ,
		memberRepo: repository.NewMemberRepository(),
		coachRepo: repository.NewCoachRepository(),
		objstore: objstore,
	}
}