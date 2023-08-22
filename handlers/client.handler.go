package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/repository"
	"net/http"
	"strconv"

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
	clients, _ := h.clientRepo.Get()
	return c.Render(http.StatusOK, "admin/clients/main", Data{
		"title": "Clients",
		"module": "Clients",
		"csrf" : csrf,
		"clients": clients,
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

func (h * ClientHandler)RenderClientUpdatePage(c echo.Context) error{
	csrf := c.Get("csrf")
	id := c.Param("id")
	clientId, _ := strconv.Atoi(id)
	client, getClientErr := h.clientRepo.GetClientById(clientId)
	if getClientErr != nil {
		logger.Error(getClientErr.Error(), zap.String("error", "getClientErr"))
	}
	return c.Render(http.StatusOK, "admin/clients/update-client", Data{
		"title": "Client | Update Profile ",
		"module": "Client Profile",
		"csrf" : csrf,
		"client": client,
		"isClientExist":  client.Id > 0,
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
func (h *ClientHandler) ResetPassword(c echo.Context) error {
	id := c.Param("id")
	clientId, convErr := strconv.Atoi(id)

	if convErr != nil {
		logger.Error(convErr.Error(), zap.String("error", "convErr"))
		return c.JSON(http.StatusBadRequest, Data{
			"status": http.StatusBadRequest,
			"message": "Unknown error occured.",
		})
	}
	diceware := &acopw.Diceware{
		Separator: "-",
		Length: 2,
		Capitalize: true,
	}
	generatedPassword, generateErr  := diceware.Generate()
	if generateErr != nil {
		logger.Error(generateErr.Error(), zap.String("error", "generateErr"))
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
			"message": "Unknown error occured.",

		})
	}
	hashedPassword, hashingErr := bcrypt.GenerateFromPassword([]byte(generatedPassword), bcrypt.DefaultCost)
	if hashingErr != nil {
		logger.Error(hashingErr.Error(), zap.String("error", "hashingErr"))
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
			"message": "Unknown error occured.",})

	}
	updateErr := h.clientRepo.UpdatePassword(string(hashedPassword), clientId)
	if updateErr!= nil {
		logger.Error(updateErr.Error(), zap.String("error", "updateErr"))
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
			"message": "Unknown error occured.",})

	}
	return c.JSON(http.StatusOK, Data{
		"status": http.StatusOK,
		"message": "Password Reset",
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