package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)


type ContactUs struct {


}
func NewContactUsHandler () ContactUs {
	return ContactUs{}
}

func(h * ContactUs) RenderContactUs(c echo.Context) error{ 
	return c.Render(http.StatusOK, "public/contactus", Data{})
}