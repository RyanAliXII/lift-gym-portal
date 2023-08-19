package handlers

import "github.com/labstack/echo/v4"

func RegisterHandlers(router *echo.Echo) {
	loginHandler := NewLoginHandler()
	dashboardHandler := NewDashboardHandler()


	
	router.GET("/login", loginHandler.RenderLoginPage)
	router.GET("/dashboard", dashboardHandler.RenderDashboardPage)
}