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

type CoachHandler struct {
	coachRepo repository.CoachRepository
}

func (h *CoachHandler) RenderCoachPage(c echo.Context) error {
	csrf := c.Get("csrf")
	coaches,_ := h.coachRepo.GetCoaches()
	return c.Render(http.StatusOK, "admin/coach/main", Data{
		"title": "Coaches",
		"module": "Coaches",
		"csrf": csrf,
		"coaches": coaches,
	})
}
func (h * CoachHandler)RenderCoachRegistrationPage(c echo.Context) error {
	csrf := c.Get("csrf")
	return c.Render(http.StatusOK, "admin/coach/register-coach", Data{
		"csrf": csrf,
		"title": "Coach | Registration",
		"module": "Registration Form",
	} )
}


func (h  *CoachHandler)NewCoach (c echo.Context) error {
	coach := model.Coach{}
	c.Bind(&coach)
	err, fieldErrs := coach.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": fieldErrs,
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
	coach.Password = string(hashedPassword)
	if hashingErr != nil {
		logger.Error(hashingErr.Error(), zap.String("error", "hashingErr"))
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
			"message": "Unknown error occured.",})

	}
	err = h.coachRepo.NewCoach(coach)
	if err != nil {	
		logger.Error(err.Error(), zap.String("error", "newCoach"))
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
			"message": "Unknown error occured.",})
	}
	
	return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: Data{
				"password": generatedPassword,
			},
			Message: "Coach has been registered.",
	})
}

func NewCoachHandler() CoachHandler{
	return CoachHandler{
		coachRepo: repository.NewCoachRepository(),
	}
}