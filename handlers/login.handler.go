package handlers

import (
	"fmt"
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/mysqlsession"
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type LoginHandler struct {
	userRepository  repository.UserRepository
}



func (h *LoginHandler) RenderLoginPage(c echo.Context) error{
	return c.Render(http.StatusOK, "admin/login/main", Data{
		"title": "Sign In",
	} )
}

func (h * LoginHandler) Login (c echo.Context) error {
	user := model.User{}
	bindErr := c.Bind(&user)
	if bindErr != nil {
		return c.JSON(http.StatusBadRequest, Data{
			 "status": http.StatusBadRequest,
			"message": "Invalid email or password.",
		})
	}
	dbUser, getUserErr  := h.userRepository.GetUserByEmail(user.Email)
	if getUserErr != nil {
		return c.JSON(http.StatusBadRequest, Data{
			"status": http.StatusBadRequest,
		   "message": "Invalid email or password.",
	   })
	}
	
	comparePassErr := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if comparePassErr != nil {
		return c.JSON(http.StatusBadRequest, Data{
			"status": http.StatusBadRequest,
		   "message": "Invalid email or password.",
	   })
	}
	s , getSessionErr := session.Get("sid", c)
	if getSessionErr != nil {
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
		   "message": "Unknown error occured",
	   })
	}
	s.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 24, // 1 day
		HttpOnly: true,
	}
	s.Values["data"] = mysqlsession.SessionData{
		User: mysqlsession.SessionUser{
			Id: dbUser.Id,
			GivenName: dbUser.GivenName,
			MiddleName: dbUser.MiddleName,
			Surname: dbUser.Surname,
			Email: dbUser.Email,
		},
	}
	saveErr := s.Save(c.Request(), c.Response())
	if saveErr != nil {
		fmt.Println(saveErr)
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
		   "message": "Unknown error occured",
	   })
	}
	return c.JSON(http.StatusOK, Data{
		"status": http.StatusOK,
	   "message": "Success.",
   })
}
func NewLoginHandler() LoginHandler{
	return LoginHandler{
		userRepository:  repository.NewUserRepository(),
	}
}




