package routes

import (
	"taas-api/handlers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRoutes initializes all API routes
func SetupRoutes() *gin.Engine {
	router := gin.Default()
	// CORS middleware configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins, change to specific origins as needed
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,           // Allow cookies and credentials
		MaxAge:           12 * time.Hour, // Cache preflight requests for 12 hours
	}))

	// Authentication User routes
	router.POST("/api/signup", handlers.Signup)
	router.POST("/api/login", handlers.Login)

	router.POST("/api/register", handlers.RegisterUser)
	router.POST("/api/signin", handlers.SignIn)
	router.POST("/api/register-talent", handlers.RegisterTalent)
	router.GET("/api/get-talents", handlers.GetTalentAccounts)

	// Card routes
	router.POST("/api/cards", handlers.SaveCard)
	router.POST("/api/generate-cards", handlers.CreateServiceCard)
	router.GET("/api/get-cards", handlers.GetCardsByTalentID)
	router.GET("/api/cards/user", handlers.GetUserCards)
	router.GET("/api/cards/all", handlers.GetAllCards)
	router.GET("/api/cards/event-id", handlers.Cards_Id)

	//Video Control routes
	router.POST("/api/save-video-control", handlers.SaveVideoControl)
	router.GET("/api/get-video-control", handlers.GetVideoControl)
	router.POST("/api/update-video-control", handlers.UpdateVideoControl)

	//Booking routes
	router.GET("/api/bookings/user/:user_id", handlers.GetBookingsByUser)
	router.GET("/api/bookings/talent/:talent_id", handlers.GetBookingsByTalent)
	router.POST("/api/book-cards", handlers.CreateBookings)

	router.GET("/api/bookingRequest", handlers.RetrieveMyBookedCardsRequestToTalent)
	router.PATCH("/api/handle-bookingStatus", handlers.HandleUpdateBookingStatus)

	//Notification route
	router.GET("/api/notifications", handlers.HandleNotificationStream)

	//Static files route
	router.POST("/api/upload", handlers.FileUploadHandler)

	//Schedule management routes
	router.POST("/api/create-schedule", handlers.CreateAvailableSlots)
	router.GET("/api/get-all-schedule", handlers.GetTalentAvailability)
	// router.PATCH("/api/update-schedule", handlers.UpdateAvailableSlots)
	// router.DELETE("/api/delete-schedule", handlers.DeleteAvailableSlot)
	return router
}
