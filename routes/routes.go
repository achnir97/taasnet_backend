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

	// Card routes
	router.POST("/api/cards", handlers.SaveCard)
	router.GET("/api/cards/user", handlers.GetUserCards)
	router.GET("/api/cards/all", handlers.GetAllCards)
	router.GET("/api/cards/event-id", handlers.Cards_Id)

	//Booking routes
	router.POST("/api/bookings", handlers.BookCard)
	router.GET("/api/mybookings", handlers.RetrieveMyBookedCards)

	//Video Control routes
	router.POST("/api/save-video-control", handlers.SaveVideoControl)
	router.GET("/api/get-video-control", handlers.GetVideoControl)
	router.POST("/api/update-video-control", handlers.UpdateVideoControl)
	router.GET("/api/bookingRequest", handlers.RetrieveMyBookedCardsRequestToTalent)
	router.PATCH("/api/handle-bookingStatus", handlers.HandleUpdateBookingStatus)
	router.GET("/api/notifications", handlers.HandleNotificationStream)
	router.POST("api/upload", handlers.FileUploadHandler)

	return router

}
