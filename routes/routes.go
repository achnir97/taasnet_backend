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
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,           // Allow cookies and credentials
		MaxAge:           12 * time.Hour, // Cache preflight requests for 12 hours
	}))

	// User routes
	router.POST("/api/signup", handlers.Signup)
	router.POST("/api/login", handlers.Login)
	// Card routes

	router.POST("/api/cards", handlers.SaveCard)
	//Booking routes

	router.POST("/api/bookings", handlers.BookEvent)
	//Video Control routes
	router.GET("/api/cards/user", handlers.GetUserCards)
	router.GET("/api/cards/event-id", handlers.Cards_Id)

	router.POST("/api/save-video-control", handlers.SaveVideoControl)
	router.GET("/api/get-video-control", handlers.GetVideoControl)
	router.POST("/api/update-video-control", handlers.UpdateVideoControl)

	return router

}
