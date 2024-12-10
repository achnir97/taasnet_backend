package routes

import (
	"taas-api/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes initializes all API routes
func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// User routes
	router.POST("/signup", handlers.Signup)
	router.POST("/login", handlers.Login)

	// Card routes
	router.POST("/cards", handlers.SaveCard)
	// Booking routes
	router.POST("/bookings", handlers.BookEvent)
	// Video Control routes
	router.POST("/video-control", handlers.SaveVideoControl)
	return router
}
