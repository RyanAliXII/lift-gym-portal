package handlers

import (
	"lift-fitness-gym/app/repository"
	"net/http"

	"github.com/labstack/echo/v4"
)


type MessageHandler struct {
	messsageRepo  repository.Message
}

func NewMessageHandler() MessageHandler {
	return MessageHandler{
		messsageRepo: repository.NewMessageRepository(),
	}
}
func (h * MessageHandler)RenderPage(c echo.Context) error {
	contentType := c.Request().Header.Get("content-type")
	if contentType == "application/json" {
		messages, err := h.messsageRepo.GetMessages()
		if err != nil {
			logger.Error(err.Error())
		}
		return c.JSON(http.StatusOK, JSONResponse{
			Status: http.StatusOK,
			Data: Data{
				"messages":  messages,
			},
		} )
	}
	return c.Render(http.StatusOK, "admin/messages/main", Data{ "title": "Messages", "module": "Messages"})
}