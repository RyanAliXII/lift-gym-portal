package handlers

import (
	"lift-fitness-gym/app/http/middlewares"

	"github.com/labstack/echo/v4"
)

func RegisterHandlers(router *echo.Echo) {
	

	adminRoutes(router.Group("/app"))
	clientRoutes(router.Group("/clients"))
}

func adminRoutes (router  * echo.Group){
	dashboardHandler := NewDashboardHandler()
	packageHandler := NewPackageHandler()
	clientHandler := NewClientHandler()
	membersHandler := NewMembersHandler()
	membershipPlanHandler := NewMembershipPlanHandler()
	loginHandler := NewLoginHandler()
	coachHandler := NewCoachHandler()
	router.GET("/login", loginHandler.RenderAdminLoginPage)
	router.POST("/login", loginHandler.Login)
	router.Use(middlewares.AuthMiddleware)
	router.GET("/dashboard", dashboardHandler.RenderDashboardPage,)
	router.GET("/packages", packageHandler.RenderPackagePage)
	router.POST("/packages", packageHandler.NewPackage)
	router.PUT("/packages/:id", packageHandler.UpdatePackage)
	router.GET("/clients", clientHandler.RenderClientPage)
	router.GET("/clients/:id", clientHandler.RenderClientUpdatePage)
	router.PUT("/clients/:id", clientHandler.UpdateClient)
	router.POST("/clients", clientHandler.NewClient)
	router.PATCH("/clients/:id/password", clientHandler.ResetPassword)
	router.GET("/clients/registration", clientHandler.RenderClientRegistrationForm)
	router.GET("/members", membersHandler.RenderMembersPage)
	router.POST("/members", membersHandler.Subscribe)
	router.DELETE("/subscriptions/:subscriptionId", membersHandler.CancelSubscription)
	router.GET("/memberships", membershipPlanHandler.RenderMembershipPlanPage)
	router.POST("/memberships", membershipPlanHandler.NewMembershipPlan)
	router.PUT("/memberships/:id", membershipPlanHandler.UpdatePlan)
	router.GET("/coaches", coachHandler.RenderCoachPage)
	router.GET("/coaches/registration", coachHandler.RenderCoachRegistrationPage)
	router.GET("/coaches/:id", coachHandler.RenderCoachUpdatePage)
	router.POST("/coaches", coachHandler.NewCoach)
	router.PUT("/coaches/:id", coachHandler.UpdateCoach)
	router.PATCH("/coaches/:id/password", coachHandler.ResetPassword)
	
}


func clientRoutes(router * echo.Group){
	loginHandler := NewLoginHandler()
	dashboardHandler := NewDashboardHandler()
	profileHandler := NewProfileHandler()
	verificationHandler := NewVerificationHandler()
	membershipRequestHandler := NewMembershipRequestHandler()
	router.GET("/login", loginHandler.RenderClientLoginPage)
	router.POST("/login", loginHandler.LoginClient)
	router.GET("/verification/:id",  verificationHandler.VerifyEmail)
	router.Use(middlewares.ClientAuthMiddleware)
	router.GET("/dashboard", dashboardHandler.RenderClientDashboard)
	router.GET("/profile", profileHandler.RenderClientProfilePage)
	router.POST("/profile/verification", profileHandler.CreateEmailVerification)
	router.PATCH("/profile/password", profileHandler.ChangePassword)
	router.GET("/membership-requests", membershipRequestHandler.RenderClientMembershipRequest)
}