package handlers

import (
	"fmt"
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/mailer"
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/app/repository"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type ProfileHandler struct {
	clientRepo repository.ClientRepository
	verificationRepo  repository.VerificationRepository
	memberRepo repository.MemberRepository
	coachRepo repository.CoachRepository

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
func NewProfileHandler()ProfileHandler {
	return ProfileHandler{
		clientRepo: repository.NewClientRepository(),
		verificationRepo: repository.NewVerificationRepository() ,
		memberRepo: repository.NewMemberRepository(),
		coachRepo: repository.NewCoachRepository(),
	}
}