package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/repository"
	"net/http"

	"git.sr.ht/~jamesponddotco/acopw-go"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type ClientHandler struct {
	clientRepo repository.ClientRepository
}
func (h * ClientHandler) RenderClientPage(c echo.Context) error {
	csrf := c.Get("csrf")
	return c.Render(http.StatusOK, "admin/clients/main", Data{
		"title": "Clients",
		"module": "Clients",
		"csrf" : csrf,
	})
}
func (h * ClientHandler) RenderClientRegistrationForm(c echo.Context) error {
	csrf := c.Get("csrf")
	return c.Render(http.StatusOK, "admin/clients/register-client", Data{
		"title": "Client | Registration ",
		"module": "Registration Form",
		"csrf" : csrf,
	})
}
func (h *ClientHandler) NewClient(c echo.Context) error {
	client := model.Client{}
	bindErr := c.Bind(&client)
	if bindErr != nil {
		logger.Error(bindErr.Error(), zap.String("error", "bindErr"))
		return c.JSON(http.StatusBadRequest, Data{
			"status": http.StatusBadRequest,
			"message": "Unknown error occured.",

		})
	}
	validateErr, fieldErrs := client.Validate()
	if validateErr != nil {
		logger.Error(validateErr.Error(), zap.String("error", "validateErr"))
		return c.JSON(http.StatusBadRequest, Data{
			"status": http.StatusBadRequest,
			"message": "Validation error.",
			"data" : Data{
				"errors" : fieldErrs,
			},

		})
	}
	diceware := &acopw.Diceware{
		Separator: "-",
		Length: 2,
		Capitalize: true,
	}
	generatedPassword,generateErr  := diceware.Generate()

	if generateErr != nil {
		logger.Error(generateErr.Error(), zap.String("error", "generateErr"))
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
			"message": "Unknown error occured.",

		})
	}
	hashedPassword, hashingErr := bcrypt.GenerateFromPassword([]byte(generatedPassword), bcrypt.DefaultCost)
	client.Password = string(hashedPassword)
	if hashingErr != nil {
		logger.Error(hashingErr.Error(), zap.String("error", "hashingErr"))
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
			"message": "Unknown error occured.",})

	}
	newClientErr := h.clientRepo.New(client)
	if newClientErr != nil {
		logger.Error(newClientErr.Error(), zap.String("error", "newClientErr"))
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
			"message": "Unknown error occured.",

		})
	}
	return c.JSON(http.StatusOK, Data{
		"status": http.StatusOK,
		"message": "Client registered.",
		"data": Data{
			"password": generatedPassword,
		},
	})
}


func NewClientHandler() ClientHandler {
	return ClientHandler{
		clientRepo: repository.NewClientRepository(),
	}
}