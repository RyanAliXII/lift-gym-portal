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
	packageHandlder := NewPackageHandler()
	router.Use(middlewares.AuthMiddleware)
	router.GET("/dashboard", dashboardHandler.RenderDashboardPage,)
	router.GET("/packages", packageHandlder.RenderPackagePage)
	router.POST("/packages", packageHandlder.NewPackage)
	router.PUT("/packages/:id", packageHandlder.UpdatePackage)
}