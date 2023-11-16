package handlers

import (
	"lift-fitness-gym/app/model"
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)


type ContactUs struct {

	messageRepo repository.Message
}
func NewContactUsHandler () ContactUs {
	return ContactUs{
		messageRepo: repository.NewMessageRepository(),
	}
}

func(h * ContactUs) RenderContactUs(c echo.Context) error{ 
	return c.Render(http.StatusOK, "public/contactus", Data{})
}
func (h * ContactUs) NewMessage (c echo.Context) error {
	message := model.ContactUs{}

	c.Bind(&message)
	fields, err := message.Validate()
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "ValidateError"))
		return c.JSON(http.StatusBadRequest, JSONResponse{
			Status: http.StatusBadRequest,
			Data: Data{
				"errors": fields,
			},
			Message: "Validation error",
		})
	}
	err = h.messageRepo.NewMessage(message)
	if err != nil {
		logger.Error(err.Error(), zap.String("error", "NewMessageErr"))
		return c.JSON(http.StatusInternalServerError, JSONResponse{
			Status: http.StatusInternalServerError,
			Message: "Unknown error occured.",
		})
	}
	return c.JSON(http.StatusOK, JSONResponse{
		Status: http.StatusOK,
		Message: "Message sent.",
	})
}
