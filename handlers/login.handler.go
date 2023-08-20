package handlers

import (
	"fmt"
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo/v4"
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
			"message": "Invalid username or password.",
		})
	}
	fetchedUser, getUserErr  := h.userRepository.GetUserByEmail(user.Email)
	if getUserErr != nil {
		return c.JSON(http.StatusBadRequest, Data{
			"status": http.StatusBadRequest,
		   "message": "Invalid username or password.",
	   })
	}
	fmt.Println(fetchedUser)
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




