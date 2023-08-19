package handlers

import "github.com/labstack/echo/v4"

func RegisterHandlers(router *echo.Echo) {
	loginHandler := NewLoginHandler()
	dashboardHandler := NewDashboardHandler()
	packageHandlder := NewPackageHandler()

	
	router.GET("/login", loginHandler.RenderLoginPage)
	router.GET("/dashboard", dashboardHandler.RenderDashboardPage)
	router.GET("/packages", packageHandlder.RenderPackagePage)
	router.POST("/packages", packageHandlder.NewPackage)
}