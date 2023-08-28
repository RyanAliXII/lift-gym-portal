package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/mailer"
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ProfileHandler struct {
	clientRepo repository.ClientRepository
	verificationRepo  repository.VerificationRepository
}

func (h *ProfileHandler) RenderClientProfilePage(c echo.Context) error{
	csrf := c.Get("csrf")
	s, getSessionErr := session.Get("client_sid", c)
	if getSessionErr != nil {
		logger.Error(getSessionErr.Error(), zap.String("error",  "getSessionErr"))
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
		   "message": "Unknown error occured",
	   })
	}

	sessionData := mysqlsession.SessionData{}
	bindErr := sessionData.Bind(s.Values["data"])
	if bindErr != nil {
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
		   "message": "Unknown error occured",
	   })
	}
	client, getClientErr := h.clientRepo.GetById(sessionData.User.Id)
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
	})
	return nil
}
func (h * ProfileHandler) CreateEmailVerification(c echo.Context) error {
	s, getSessionErr := session.Get("client_sid", c)
	if getSessionErr != nil {
		logger.Error(getSessionErr.Error(), zap.String("error",  "getSessionErr"))
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
		   "message": "Unknown error occured",
	   })
	}
	sessionData := mysqlsession.SessionData{}
    bindErr := sessionData.Bind(s.Values["data"])
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

func NewProfileHandler()ProfileHandler {
	return ProfileHandler{
		clientRepo: repository.NewClientRepository(),
		verificationRepo: repository.NewVerificationRepository() ,
	}
}