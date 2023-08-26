package handlers

import (
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

type ProfileHandler struct {
	clientRepo repository.ClientRepository
}

func (h *ProfileHandler) RenderClientProfilePage(c echo.Context) error{
	csrf := c.Get("csrf")
	s, getSessionErr := session.Get("sid", c)
	if getSessionErr != nil {
		logger.Error(getSessionErr.Error(), zap.String("error",  "getSessionErr"))
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
		   "message": "Unknown error occured",
	   })
	}
	sessionData := mysqlsession.SessionData{}
	decodeErr := mapstructure.Decode(s.Values["data"], &sessionData)
	if decodeErr != nil {
		logger.Error(decodeErr.Error(), zap.String("error",  "decodeErr"))
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
	c.Render(http.StatusOK, "client/profile/main", Data{
		"csrf": csrf,
		"title": "Profile",
		"module": "Profile",
		"profile": client,
	})
	return nil
}

func NewProfileHandler()ProfileHandler {
	return ProfileHandler{
		clientRepo: repository.NewClientRepository(),
	}
}