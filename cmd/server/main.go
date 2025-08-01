package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/uso/uso/config"
	"github.com/uso/uso/internal/database"
	"github.com/uso/uso/internal/handlers"
	"github.com/uso/uso/internal/middleware"
	"github.com/stripe/stripe-go/v76"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize Stripe
	stripe.Key = cfg.StripeSecretKey

	// Connect to database
	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migrations check
	if err := db.RunMigrations(); err != nil {
		log.Printf("Migration check failed: %v", err)
	}

	// Initialize Gin router
	router := gin.Default()
	router.Use(middleware.CORS())

	// Static files
	router.Static("/static", "./web/static")
	router.LoadHTMLGlob("web/templates/*")

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, cfg)
	castHandler := handlers.NewCastHandler(db, cfg)
	bookingHandler := handlers.NewBookingHandler(db, cfg)
	searchHandler := handlers.NewSearchHandler(db)
	adminHandler := handlers.NewAdminHandler(db, cfg)

	// Public routes
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// API routes
	api := router.Group("/api")
	{
		// Auth routes
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)

		// Search routes (public)
		api.GET("/casts/search", searchHandler.SearchCasts)
		api.GET("/casts/:id", searchHandler.GetCastProfile)
		api.GET("/service-areas", searchHandler.GetServiceAreas)

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthRequired(cfg))
		{
			// User profile
			protected.GET("/profile", authHandler.GetProfile)
			protected.PUT("/profile", authHandler.UpdateProfile)
			protected.POST("/profile/image", authHandler.UploadProfileImage)

			// Cast routes
			castRoutes := protected.Group("/cast")
			castRoutes.Use(middleware.CastOnly())
			{
				castRoutes.PUT("/profile", castHandler.UpdateCastProfile)
				castRoutes.POST("/gallery", castHandler.UploadGalleryImage)
				castRoutes.DELETE("/gallery/:id", castHandler.DeleteGalleryImage)
				castRoutes.GET("/bookings", castHandler.GetCastBookings)
				castRoutes.POST("/bookings/:id/respond", castHandler.RespondToBooking)
				castRoutes.GET("/earnings", castHandler.GetEarnings)
			}

			// Guest routes
			guestRoutes := protected.Group("/guest")
			guestRoutes.Use(middleware.GuestOnly())
			{
				guestRoutes.POST("/bookings", bookingHandler.CreateBooking)
				guestRoutes.GET("/bookings", bookingHandler.GetGuestBookings)
				guestRoutes.POST("/bookings/:id/cancel", bookingHandler.CancelBooking)
			}

			// Shared booking routes
			protected.GET("/bookings/:id", bookingHandler.GetBooking)
			protected.POST("/bookings/:id/complete", bookingHandler.CompleteBooking)

			// Messages (only for accepted bookings)
			protected.GET("/bookings/:id/messages", bookingHandler.GetMessages)
			protected.POST("/bookings/:id/messages", bookingHandler.SendMessage)

			// Reviews
			protected.POST("/reviews", bookingHandler.CreateReview)
			protected.GET("/users/:id/reviews", bookingHandler.GetUserReviews)
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.AdminOnly(cfg))
		{
			admin.GET("/dashboard", adminHandler.GetDashboard)
			admin.GET("/casts/pending", adminHandler.GetPendingCasts)
			admin.POST("/casts/:id/approve", adminHandler.ApproveCast)
			admin.POST("/casts/:id/reject", adminHandler.RejectCast)
			admin.GET("/bookings", adminHandler.GetAllBookings)
			admin.GET("/analytics", adminHandler.GetAnalytics)
		}

		// Stripe webhook
		api.POST("/stripe/webhook", bookingHandler.HandleStripeWebhook)
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}