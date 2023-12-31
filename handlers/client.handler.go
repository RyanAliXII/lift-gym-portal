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
	
	contentType := c.Request().Header.Get("content-type")
	
	if contentType == "application/json" {
		return h.handleJSONContent(c)
	}
	clients, _ := h.clientRepo.Get()
	return c.Render(http.StatusOK, "admin/clients/main", Data{
		"title": "Clients",
		"module": "Clients",
		"csrf" : csrf,
		"clients": clients,
	})
}

func (h * ClientHandler) handleJSONContent( c echo.Context ) error {


	keyword := c.QueryParam("keyword")
	if len(keyword) > 0 {
		clients, err := h.clientRepo.Search(keyword)
	
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "Search"))
		}
		return c.JSON(http.StatusOK, JSONResponse{
			   Status: http.StatusOK,
			   Data: Data{
				   "clients": clients,
			   },
			   Message: "Searched Client.",
		})

	}
	clientType := c.QueryParam("type")
	if len(clientType) > 0 {
		if clientType == "unsubscribed"{
		 clients, _ := h.clientRepo.GetUnsubscribed()
		 return c.JSON(http.StatusOK, JSONResponse{
				Status: http.StatusOK,
				Data: Data{
					"clients": clients,
				},
				Message: "Unsubscibed clients fetched.",
			})
		}
	}
	clients, _ := h.clientRepo.Get()
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Data: Data{
			"clients": clients,
		},
		Message: "Clients fetched.",
	})
}
func (h * ClientHandler) DeleteClient( c echo.Context) error {
	id := c.Param("id")
	clientId, convErr := strconv.Atoi(id)
	if convErr != nil {
		logger.Error(convErr.Error(), zap.String("error", "convErr"))
		return c.JSON(http.StatusBadRequest, Data{
			"status": http.StatusBadRequest,
			"message": "Unknown error occured.",
		})
	}

	err := h.clientRepo.Delete(clientId)
	if err != nil{
		logger.Error(err.Error(), zap.String("error","deleteErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",	
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Client deleted.",
	})
}


func (h * ClientHandler) VerifyClient( c echo.Context) error {
	id := c.Param("id")
	clientId, convErr := strconv.Atoi(id)
	if convErr != nil {
		logger.Error(convErr.Error(), zap.String("error", "convErr"))
		return c.JSON(http.StatusBadRequest, Data{
			"status": http.StatusBadRequest,
			"message": "Unknown error occured.",
		})
	}
	
	err := h.clientRepo.MarkAsVerified(clientId)
	if err != nil{
		logger.Error(err.Error(), zap.String("error","verifyClient"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",	
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Client verified.",
	})
}

func (h * ClientHandler) RemoveVerification( c echo.Context) error {
	id := c.Param("id")
	clientId, convErr := strconv.Atoi(id)
	if convErr != nil {
		logger.Error(convErr.Error(), zap.String("error", "convErr"))
		return c.JSON(http.StatusBadRequest, Data{
			"status": http.StatusBadRequest,
			"message": "Unknown error occured.",
		})
	}
	
	err := h.clientRepo.MarkAsUnverified(clientId)
	if err != nil{
		logger.Error(err.Error(), zap.String("error","removeVerification"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",	
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Client unverified.",
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

func (h *ClientHandler) UpdateClient(c echo.Context) error {
	client := model.Client{}
	bindErr := c.Bind(&client)
	if bindErr != nil {
		logger.Error(bindErr.Error(), zap.String("error", "bindErr"))
		return c.JSON(http.StatusBadRequest, Data{
			"status": http.StatusBadRequest,
			"message": "Unknown error occured.",

		})
	}
	validateErr, fieldErrs := client.ValidateUpdate()
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
	updateClientErr := h.clientRepo.Update(client)
	if updateClientErr != nil {
		logger.Error(updateClientErr .Error(), zap.String("error", "updateClientErr"))
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
			"message": "Unknown error occured.",

		})
	}
	return c.JSON(http.StatusOK, Data{
		"status": http.StatusOK,
		"message": "Client updated.",
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
func (h * ClientHandler) RenderMembersPage ( c echo.Context) error {
	return c.Render(http.StatusOK,"admin/members/main", Data{})
}
func NewClientHandler() ClientHandler {
	return ClientHandler{
		clientRepo: repository.NewClientRepository(),
	}
}