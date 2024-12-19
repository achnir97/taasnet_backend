package handlers

import (
	"fmt"
	"net/http"
	"taas-api/config"
	"taas-api/models"

	"github.com/gin-gonic/gin"

	"time"
)

func HandleNotificationStream(c *gin.Context) {
	// Set SSE headers
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	// Extract user ID from query
	userId := c.Query("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Infinite loop to simulate real-time updates
	for {
		var notifications []models.Booking

		// Fetch all pending bookings for the user
		if err := config.DB.Where("user_id = ? AND status = ?", userId, "Pending").Find(&notifications).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
			return
		}

		// If there are pending notifications, send them to the client
		if len(notifications) > 0 {
			for _, booking := range notifications {
				data := fmt.Sprintf(`{"id": "%s", "eventId": "%s", "message": "Booking '%s' is pending", "status": "%s", "bookedBy": "%s"}`,
					booking.ID, booking.EventID, booking.Title, booking.Status, booking.BookedBy)
				fmt.Fprintf(c.Writer, "data: %s\n\n", data)
				c.Writer.Flush() // Push data to the client
			}
		}

		// Wait for 5 seconds before checking again
		time.Sleep(5 * time.Second)
	}
}
