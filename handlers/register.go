package handlers

import (
	"lift-fitness-gym/app/http/middlewares"

	"github.com/labstack/echo/v4"
)

func RegisterHandlers(router *echo.Echo) {
	loginHandler := NewLoginHandler()
	dashboardHandler := NewDashboardHandler()
	packageHandlder := NewPackageHandler()

	
	router.GET("/login", loginHandler.RenderLoginPage)
	router.POST("/login", loginHandler.Login)
	router.GET("/dashboard", dashboardHandler.RenderDashboardPage, middlewares.AuthMiddleware)
	router.GET("/packages", packageHandlder.RenderPackagePage)
	router.POST("/packages", packageHandlder.NewPackage)
	router.PUT("/packages/:id", packageHandlder.UpdatePackage)
}