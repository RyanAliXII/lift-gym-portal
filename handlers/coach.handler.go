package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/pkg/objstore"
	"lift-fitness-gym/app/repository"
	"net/http"
	"strconv"

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

func (h * CoachHandler) RenderClientHireCoachPage (c echo.Context ) error {
	contentType := c.Request().Header.Get("Content-Type")
	if contentType == "application/json"{
		coaches, err := h.coachRepo.GetCoaches()
		if err != nil {
			logger.Error(err.Error(), zap.String("error", "GetCoachesErr"))
		}

		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: Data{
				"coaches": coaches,
			},
			Message: "Coaches fetched.",
		})
	}
	return c.Render(http.StatusOK, "client/hire-a-coach/main", Data{
		"csrf": c.Get("csrf"),
		"title": "Hire a Coach",
		"module": "Coaches",
		"objstorePublicUrl": objstore.PublicURL,		
 	})
}

func (h * CoachHandler) HireCoach (c echo.Context ) error {
	hiredCoach := model.HiredCoach{}


	err := c.Bind(&hiredCoach)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "bindErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	err, fields := hiredCoach.Validate()

	if err != nil {

		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": fields,
			},
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Coach has been hired.",
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
func (h * CoachHandler) RenderCoachUpdatePage(c echo.Context) error {
	csrf := c.Get("csrf")
	id := c.Param("id")
	coachId, _ := strconv.Atoi(id)
	coach, getClientErr := h.coachRepo.GetCoachById(coachId)
	if getClientErr != nil {
		logger.Error(getClientErr.Error(), zap.String("error", "getClientErr"))
	}
	return c.Render(http.StatusOK, "admin/coach/update-coach", Data{
		"title": "Coach | Update Profile ",
		"module": "Coach Profile",
		"csrf" : csrf,
		"coach": coach,
		"isCoachExist":  coach.Id > 0,
	})



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

func (h  *CoachHandler)UpdateCoach(c echo.Context) error {
	coach := model.Coach{}
	c.Bind(&coach)

	err, fieldErrs := coach.ValidateUpdate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": fieldErrs,
			},
		})
	}
	
	err = h.coachRepo.UpdateCoach(coach)
	if err != nil {	
		logger.Error(err.Error(), zap.String("error", "updateCoach"))
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
			"message": "Unknown error occured.",})
	}
	
	return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: nil,
			Message: "Coach has profile updated.",
	})
}
func (h  *CoachHandler)ResetPassword(c echo.Context) error {
	id := c.Param("id")
	coachId, convErr := strconv.Atoi(id)
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
	generatedPassword,generateErr  := diceware.Generate()

	if generateErr != nil {
		logger.Error(generateErr.Error(), zap.String("error", "generateErr"))
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
			"message": "Unknown error occured.",

		})
	}
	hashedPassword, hashingErr := bcrypt.GenerateFromPassword([]byte(generatedPassword), bcrypt.DefaultCost)
	h.coachRepo.UpdatePassword(string(hashedPassword), coachId)
	if hashingErr != nil {
		logger.Error(hashingErr.Error(), zap.String("error", "hashingErr"))
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
			"message": "Unknown error occured.",})

	}
	
	return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: Data{
				"password": generatedPassword,
			},
			Message: "Password has been updated.",
	})
}


func NewCoachHandler() CoachHandler{
	return CoachHandler{
		coachRepo: repository.NewCoachRepository(),
	}
}