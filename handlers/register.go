package handlers

import (
	"lift-fitness-gym/app/http/middlewares"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterHandlers(router *echo.Echo) {
	passwordHandler := NewPasswordHandler()
	router.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			contentType := c.Request().Header.Get("Content-Type")
			if contentType == "application/json" {
				c.Response().Header().Set("Vary", "Accept")
				return next(c)
			}
			return next(c)
		}
	})
	
	router.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "public/landing", nil)
	})

	router.GET("/change-password", passwordHandler.RenderChangePasswordPage)
	router.POST("/change-password", passwordHandler.ChangePassword)
	adminRoutes(router.Group("/app"))
	clientRoutes(router.Group("/clients"))
	coachRoutes(router.Group("/coaches"))
}

func adminRoutes (router  * echo.Group){
	dashboardHandler := NewDashboardHandler()
	packageHandler := NewPackageHandler()
	clientHandler := NewClientHandler()
	membersHandler := NewMembersHandler()
	membershipPlanHandler := NewMembershipPlanHandler()
	loginHandler := NewLoginHandler()
	coachHandler := NewCoachHandler()
	pkgRequestHandler := NewPackageRequestHandler()
	membershipRequestHandler := NewMembershipRequestHandler()
	inventoryHandler := NewInventoryHandler()
	staffHandler := NewStaffHandler()
	workoutCategoryHandler := NewWorkoutCategoryHandler()
	workoutHandler := NewWorkoutHandler()
	rolesPermissionHandler := NewRoleHandler()
	clientLogHandler := NewClientLogHandler()
	passwordHandler := NewPasswordHandler()
	logoutHandler := NewLogoutHandler()
	adminProfileHandler := NewAdminProfileHandler()
	dateSlotHandler := NewDateSlotHandler()
	timeSlotHandler := NewTimeSlotHandler()
	reservationHandler := NewReservationHandler()
	PaymentHistory := NewPaymentHistoryHandler()
	reportHandler := NewReportHandler()
	coachingPenaltyHandler := NewCoachingPenalty()
	router.GET("/login", loginHandler.RenderAdminLoginPage)
	router.POST("/login", loginHandler.Login)
	router.GET("/reset-password", passwordHandler.RenderResetPasswordPage)
	router.POST("/reset-password", passwordHandler.ResetPassword)
	router.GET("/reports/:id", reportHandler.RenderReportData)
	router.Use(middlewares.AuthMiddleware("sid", "/app/login"))
	router.DELETE("/logout", logoutHandler.LogoutAdmin)
	router.GET("/dashboard", dashboardHandler.RenderDashboardPage,)
	router.GET("/packages", packageHandler.RenderPackagePage, middlewares.ValidatePermissions("Package.Read"))
	router.POST("/packages", packageHandler.NewPackage, middlewares.ValidatePermissions("Package.Create"))
	router.PUT("/packages/:id", packageHandler.UpdatePackage, middlewares.ValidatePermissions("Package.Edit"))
	router.DELETE("/packages/:id", packageHandler.DeletePackage, middlewares.ValidatePermissions("Package.Delete"))
	router.GET("/clients", clientHandler.RenderClientPage, middlewares.ValidatePermissions("Client.Read"))
	router.GET("/clients/:id", clientHandler.RenderClientUpdatePage, middlewares.ValidatePermissions("Client.Edit"))
	router.PUT("/clients/:id", clientHandler.UpdateClient, middlewares.ValidatePermissions("Client.Edit"))
	router.POST("/clients", clientHandler.NewClient, middlewares.ValidatePermissions("Client.Create"))
	router.DELETE("/clients/:id", clientHandler.DeleteClient, middlewares.ValidatePermissions("Client.Delete") )
	router.PATCH("/clients/:id/password", clientHandler.ResetPassword, middlewares.ValidatePermissions("Client.Edit"))
	router.GET("/clients/registration", clientHandler.RenderClientRegistrationForm, middlewares.ValidatePermissions("Client.Create"))
	router.GET("/members", membersHandler.RenderMembersPage, middlewares.ValidatePermissions("Member.Read"))
	router.POST("/members", membersHandler.Subscribe, middlewares.ValidatePermissions("Member.Create"))
	router.DELETE("/subscriptions/:subscriptionId", membersHandler.CancelSubscription, middlewares.ValidatePermissions("Member.Delete"))
	router.GET("/memberships", membershipPlanHandler.RenderMembershipPlanPage, middlewares.ValidatePermissions("Plan.Read"))
	router.POST("/memberships", membershipPlanHandler.NewMembershipPlan, middlewares.ValidatePermissions("Plan.Create"))
	router.PUT("/memberships/:id", membershipPlanHandler.UpdatePlan, middlewares.ValidatePermissions("Plan.Edit"))
	router.DELETE("/memberships/:id", membershipPlanHandler.DeletePlan, middlewares.ValidatePermissions("Plan.Delete"))
	router.GET("/coaches", coachHandler.RenderCoachPage, middlewares.ValidatePermissions("Coach.Read"))
	router.GET("/coaches/registration", coachHandler.RenderCoachRegistrationPage,  middlewares.ValidatePermissions("Coach.Create"))
	router.GET("/coaches/:id", coachHandler.RenderCoachUpdatePage, middlewares.ValidatePermissions("Coach.Edit"))
	router.POST("/coaches", coachHandler.NewCoach, middlewares.ValidatePermissions("Coach.Create"))
	router.PUT("/coaches/:id", coachHandler.UpdateCoach, middlewares.ValidatePermissions("Coach.Edit"))
	router.DELETE("/coaches/:id", coachHandler.DeleteCoach, middlewares.ValidatePermissions("Coach.Delete"))
	router.PATCH("/coaches/:id/password", coachHandler.ResetPassword, middlewares.ValidatePermissions("Coach.Edit"))
	router.GET("/membership-requests", membershipRequestHandler.RenderAdminMembershipRequest, middlewares.ValidatePermissions("MembershipRequest.Read"))
	router.PATCH("/membership-requests/:id/status", membershipRequestHandler.UpdateMembershipRequestStatus, middlewares.ValidatePermissions("MembershipRequest.Edit"))
	router.GET("/package-requests", pkgRequestHandler.RenderAdminPackageRequestPage, middlewares.ValidatePermissions("PackageRequest.Read"))
	router.POST("/package-requests", pkgRequestHandler.NewPackageRequestWalkIn, middlewares.ValidatePermissions("PackageRequest.Create"))
	router.PATCH("/package-requests/:id/status", pkgRequestHandler.UpdatePackageRequestStatus, middlewares.ValidatePermissions("PackageRequest.Edit"))
	router.GET("/staffs",  staffHandler.RenderStaffPage, middlewares.ValidatePermissions("Staff.Read"))
	router.POST("/staffs",staffHandler.NewStaff, middlewares.ValidatePermissions("Staff.Create"))
	router.PUT("/staffs/:id",staffHandler.UpdateStaff, middlewares.ValidatePermissions("Staff.Edit"))
	router.DELETE("/staffs/:id",staffHandler.DeleteStaff, middlewares.ValidatePermissions("Staff.Delete"))
	router.PATCH("/staffs/:id/password",staffHandler.ResetPassword,middlewares.ValidatePermissions("Staff.Edit"))
	router.GET("/inventory", inventoryHandler.RenderInventoryPage, middlewares.ValidatePermissions("Inventory.Read"))
	router.POST("/inventory", inventoryHandler.NewEquipment, middlewares.ValidatePermissions("Inventory.Create"))
	router.PUT("/inventory/:id", inventoryHandler.UpdateEquipment, middlewares.ValidatePermissions("Inventory.Edit"))
	router.DELETE("/inventory/:id", inventoryHandler.DeleteEquipment, middlewares.ValidatePermissions("Inventory.Delete"))
	
	workoutGrp := router.Group("/workouts")
	workoutGrp.GET("", workoutHandler.RenderWorkoutPage, middlewares.ValidatePermissions("Workout.Read"))
	workoutGrp.POST("", workoutHandler.NewWorkout, middlewares.ValidatePermissions("Workout.Create"))
	workoutGrp.PUT("/:id", workoutHandler.UpdateWorkout, middlewares.ValidatePermissions("Workout.Edit"))
	workoutGrp.DELETE("/:id", workoutHandler.DeleteWorkout, middlewares.ValidatePermissions("Workout.Delete"))
	workoutGrp.GET("/categories", workoutCategoryHandler.RenderCategoryPage, middlewares.ValidatePermissions("WorkoutCategory.Read"))
	workoutGrp.POST("/categories", workoutCategoryHandler.NewCategory, middlewares.ValidatePermissions("WorkoutCategory.Create"))
	workoutGrp.PUT("/categories/:id", workoutCategoryHandler.UpdateCategory, middlewares.ValidatePermissions("WorkoutCategory.Edit"))
	workoutGrp.DELETE("/categories/:id", workoutCategoryHandler.DeleteCategory, middlewares.ValidatePermissions("WorkoutCategory.Delete"))
	router.GET("/roles", rolesPermissionHandler.RenderRolePage, middlewares.ValidatePermissions("Role.Read"))
	router.POST("/roles", rolesPermissionHandler.NewRole, middlewares.ValidatePermissions("Role.Create"))
	router.DELETE("/roles/:id", rolesPermissionHandler.DeleteRole, middlewares.ValidatePermissions("Role.Delete"))
	router.PUT("/roles/:id", rolesPermissionHandler.UpdateRole, middlewares.ValidatePermissions("Role.Edit"))
	router.GET("/client-logs", clientLogHandler.RenderClientLogPage, middlewares.ValidatePermissions("ClientLog.Read"))
	router.POST("/client-logs", clientLogHandler.NewLog, middlewares.ValidatePermissions("ClientLog.Create"))
	router.PUT("/client-logs/:id", clientLogHandler.UpdateLog, middlewares.ValidatePermissions("ClientLog.Edit"))
	router.DELETE("/client-logs/:id", clientLogHandler.DeleteLog, middlewares.ValidatePermissions("ClientLog.Delete"))
	router.GET("/profile", adminProfileHandler.RenderAdminProfile)
	router.PATCH("/profile/password", adminProfileHandler.ChangePassword)
	router.GET("/date-slots", dateSlotHandler.RenderDateSlotPage, middlewares.ValidatePermissions("DateSlot.Read"))
	router.POST("/date-slots", dateSlotHandler.NewSlot, middlewares.ValidatePermissions("DateSlot.Create"))
	router.DELETE("/date-slots/:id", dateSlotHandler.DeleteSlot,  middlewares.ValidatePermissions("DateSlot.Delete"))
	router.GET("/time-slots", timeSlotHandler.RenderTimeSlotPage, middlewares.ValidatePermissions("TimeSlot.Read"))
	router.POST("/time-slots", timeSlotHandler.NewTimeSlot, middlewares.ValidatePermissions("TimeSlot.Create"))
	router.PUT("/time-slots/:id", timeSlotHandler.UpdateTimeSlot, middlewares.ValidatePermissions("TimeSlot.Edit"))
	router.DELETE("/time-slots/:id", timeSlotHandler.DeleteTimeSlot, middlewares.ValidatePermissions("TimeSlot.Delete"))
	router.GET("/time-slots/:id/selections", timeSlotHandler.GetTimeSlotExcept, middlewares.ValidatePermissions("TimeSlot.Read"))
	router.GET("/reservations", reservationHandler.RenderAdminReservationPage, middlewares.ValidatePermissions("Reservation.Read"))
	router.GET("/reservations/date-slots/:dateSlotId", reservationHandler.GetReservationByDateSlot, middlewares.ValidatePermissions("Reservation.Read"))
	router.PUT("/reservations/:id/status", reservationHandler.UpdateReservationStatus, middlewares.ValidatePermissions("Reservation.Edit"))
	router.GET("/payments", PaymentHistory.RenderPayments, middlewares.ValidatePermissions("Payment.Read"))
	router.GET("/reports", reportHandler.RenderReportPage, middlewares.ValidatePermissions("Report.Read"))
	router.POST("/reports", reportHandler.CreateReport, middlewares.ValidatePermissions("Report.Create"))
	router.GET("/coaching-penalty", coachingPenaltyHandler.RenderPenaltyPage,middlewares.ValidatePermissions("Penalty.Read"))
	router.PATCH("/coaching-penalty/:id/settlement", coachingPenaltyHandler.SettlePenalty, middlewares.ValidatePermissions("Penalty.Edit"))
	router.PATCH("/coaching-penalty/:id/unsettlement", coachingPenaltyHandler.UnsettlePenalty, middlewares.ValidatePermissions("Penalty.Edit"))

}
func clientRoutes(router * echo.Group){
	loginHandler := NewLoginHandler()
	dashboardHandler := NewDashboardHandler()
	profileHandler := NewProfileHandler()
	verificationHandler := NewVerificationHandler()
	membershipRequestHandler := NewMembershipRequestHandler()
	pkgRequestHandler := NewPackageRequestHandler()
	workoutCategoryHandler := NewWorkoutCategoryHandler()
	registrationHandler := NewRegistrationHandler()
	paymentHistoryHandler := NewPaymentHistoryHandler()
	coachHandler := NewCoachHandler()
	coachRateHandler := NewCoachRateHandler()
	hiredCoachHandler := NewHiredCoachHandler()
	passwordHandler := NewPasswordHandler()
	logoutHandler := NewLogoutHandler()
	reservationHandler := NewReservationHandler()
	router.GET("/reset-password", passwordHandler.RenderResetClientPasswordPage)
	router.POST("/reset-password", passwordHandler.ResetClientPassword)
	router.GET("/login", loginHandler.RenderClientLoginPage)
	router.POST("/login", loginHandler.LoginClient)
	router.GET("/verification/:id",  verificationHandler.VerifyEmail)
	router.GET("/registration", registrationHandler.RenderRegistrationPage)
	router.POST("/registration", registrationHandler.RegisterClient)
	router.Use(middlewares.AuthMiddleware("client_sid", "/clients/login"))
	router.DELETE("/logout", logoutHandler.LogoutClient)
	router.GET("/dashboard", dashboardHandler.RenderClientDashboard)
	router.GET("/profile", profileHandler.RenderClientProfilePage)
	router.POST("/profile/verification", profileHandler.CreateEmailVerification)
	router.PATCH("/profile/password", profileHandler.ChangePassword)
	router.PATCH("/profile", profileHandler.UpdateProfile)
	router.GET("/membership-requests", membershipRequestHandler.RenderClientMembershipRequest)
	router.PATCH("/membership-requests/:id/status", membershipRequestHandler.CancelMembershipRequestStatus)
	router.POST("/membership-requests", membershipRequestHandler.NewRequest)
	router.GET("/memberships", membershipRequestHandler.GetUnrequestedMembershipPlans)
	router.GET("/package-requests", pkgRequestHandler.RenderClientPackageRequestPage)
	router.GET("/packages", pkgRequestHandler.GetUnrequestedPackages)
	router.POST("/package-requests", pkgRequestHandler.NewPackageRequest)
	router.PATCH("/package-requests/:id/status", pkgRequestHandler.CancelPackageRequest)
	router.GET("/workouts", workoutCategoryHandler.RenderClientWorkoutPage)
	router.GET("/workouts/:id", workoutCategoryHandler.RenderClientWorkoutsByCategoryId)
	router.GET("/hire-a-coach", coachHandler.RenderClientHireCoachPage)
	router.POST("/hire-a-coach", coachHandler.HireCoach)
	router.GET("/coaches/:coachId/rates", coachRateHandler.GetCoachRatesByCoachId)
	router.GET("/hired-coaches", hiredCoachHandler.RenderClientHiredCoachesPage)
	router.DELETE("/hired-coaches/:id", hiredCoachHandler.CancelAppointmentByClient)
	router.GET("/reservations", reservationHandler.RenderClientReservationPage)
	router.POST("/reservations", reservationHandler.NewReservation)
	router.GET("/reservations/date-slots", reservationHandler.GetDateSlots)
	router.GET("/reservations/date-slots/:dateSlotId/time-slots", reservationHandler.GetTimeSlotsBasedOnDateSlotId)
	router.PUT("/reservations/:id/status/cancellation", reservationHandler.CancelReservation)
	router.GET("/payment-history", paymentHistoryHandler.RenderClientPaymentHistory)
}

func coachRoutes(router * echo.Group) {
	loginHandler := NewLoginHandler()
	dashboardHandler := NewDashboardHandler()
	coachProfileHandler :=  NewCoachProfileHandler()
	coachRateHandler := NewCoachRateHandler()
	hiredCoachHandler := NewHiredCoachHandler()
	passwordHandler := NewPasswordHandler()
	logoutHandler := NewLogoutHandler()
	paymentHandler := NewPaymentHistoryHandler()
	workoutCategoryHandler := NewWorkoutCategoryHandler()
	workoutHandler := NewWorkoutHandler()
	router.GET("/login", loginHandler.RenderCoachLoginPage)
	router.POST("/login", loginHandler.LoginCoach)
	router.GET("/reset-password", passwordHandler.RenderResetCoachPasswordPage)
	router.POST("/reset-password", passwordHandler.ResetCoachPassword)
	router.Use(middlewares.AuthMiddleware("coach_sid", "/coaches/login"))
	router.DELETE("/logout", logoutHandler.LogoutCoach)
	router.GET("/dashboard", dashboardHandler.RenderCoachDashboard)
	router.GET("/profile", coachProfileHandler.RenderProfilePage)
	router.PATCH("/profile/password", coachProfileHandler.ChangePassword)
	router.POST("/public-profile", coachProfileHandler.UpdatePublicProfile)
	router.GET("/rates", coachRateHandler.RenderCoachRatePage)
	router.POST("/rates", coachRateHandler.NewRate)
	router.PUT("/rates/:id", coachRateHandler.UpdateRate)
	router.DELETE("/rates/:id", coachRateHandler.DeleteRate)
	router.GET("/appointments", hiredCoachHandler.RenderCoachAppointments)
	router.PATCH("/appointments/:id/status", hiredCoachHandler.UpdateStatus)
	router.GET("/payments", paymentHandler.RenderCoachPayments)
	router.GET("/workouts/categories", workoutCategoryHandler.RenderCoachCategoryPage)
	router.POST("/workouts/categories", workoutCategoryHandler.NewCategory)
	router.PUT("/workouts/categories/:id", workoutCategoryHandler.UpdateCategory)
	router.DELETE("/workouts/categories/:id", workoutCategoryHandler.DeleteCategory)
	router.GET("/workouts", workoutHandler.RenderCoachWorkoutPage)
	router.POST("/workouts", workoutHandler.NewWorkout)
	router.PUT("/workouts/:id", workoutHandler.UpdateWorkout)
	router.DELETE("/workouts/:id", workoutHandler.DeleteWorkout)
}