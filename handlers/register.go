package handlers

import (
	"lift-fitness-gym/app/http/middlewares"

	"github.com/labstack/echo/v4"
)

func RegisterHandlers(router *echo.Echo) {
	loginHandler := NewLoginHandler()
	router.GET("/login", loginHandler.RenderLoginPage)
	router.POST("/login", loginHandler.Login)
	adminRoutes(router.Group(""))
}

func adminRoutes (router  * echo.Group){
	dashboardHandler := NewDashboardHandler()
	packageHandler := NewPackageHandler()
	clientHandler := NewClientHandler()
	router.Use(middlewares.AuthMiddleware)
	router.GET("/dashboard", dashboardHandler.RenderDashboardPage,)
	router.GET("/packages", packageHandler.RenderPackagePage)
	router.POST("/packages", packageHandler.NewPackage)
	router.PUT("/packages/:id", packageHandler.UpdatePackage)
	router.GET("/clients", clientHandler.RenderClientPage)
	router.GET("/clients/:id", clientHandler.RenderClientUpdatePage)
	router.PUT("/clients/:id", clientHandler.UpdateClient)
	router.POST("/clients", clientHandler.NewClient)
	router.PATCH("clients/:id/password", clientHandler.ResetPassword)
	router.GET("/clients/registration", clientHandler.RenderClientRegistrationForm)
}