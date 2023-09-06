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

type StaffHandler struct {
	staffRepo repository.StaffRepository
}

func (h *StaffHandler) RenderStaffPage(c echo.Context)error {
	return c.Render(http.StatusOK, "admin/staff/main", Data{
		"csrf": c.Get("csrf"),
		"title": "Staffs",
		"module": "Staffs",
	})
}
func (h *StaffHandler)NewStaff (c echo.Context) error {
	staff := model.Staff{}
	err := c.Bind(&staff)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "bindErr"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
		})
	}
	err, fieldErrors := staff.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Message: "Unknown error occured.",
			Data: Data{
				"errors": fieldErrors,
			},
		})
	}
	diceware := &acopw.Diceware{
		Separator: "-",
		Length: 2,
		Capitalize: true,
	}
	generatedPassword, err := diceware.Generate()

	if err != nil {
		logger.Error(err.Error(), zap.String("error", "generateErr"))
		return c.JSON(http.StatusInternalServerError, Data{
			"status": http.StatusInternalServerError,
			"message": "Unknown error occured.",

		})
	}
	hashedPassword, err:= bcrypt.GenerateFromPassword([]byte(generatedPassword), bcrypt.DefaultCost)
	staff.Password = string(hashedPassword)
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "New staff added.",
		Data:Data{
			"password": generatedPassword,
		},
	})
}

func NewStaffHandler() StaffHandler{
	return StaffHandler{
		staffRepo: repository.NewStaffRepository(),
	}
}