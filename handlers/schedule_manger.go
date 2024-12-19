package handlers

import (
	"fmt"
	"net/http"
	"taas-api/config"
	_ "taas-api/config"
	"taas-api/models"

	"github.com/gin-gonic/gin"
)

func CreateAvailableSlots(c *gin.Context) {
	var input models.AvailableTimeSlots
	var existingAvailableSlots models.AvailableTimeSlots
	// Parse the input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("Error binding JSON:", err) // Log the error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("Received JSON payload:", input) // Debugging: Log received data
	// Check if a record exists
	if err := config.DB.First(&existingAvailableSlots).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No existing video control record found"})
		return
	}

	//Save the input from the payload
	if err := config.DB.Save(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Save New Available Slots"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "New Availiable slotes saved successfully",
	})
}

func UpdateAvailableSlots(c *gin.Context) {
	var input models.AvailableTimeSlots
	var existingAvailableSlots models.AvailableTimeSlots

	// Parse the input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("Error binding JSON:", err) // Log the error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Received JSON payload:", input) // Debugging: Log received data

	// Check if a record exists for the given ID
	if err := config.DB.First(&existingAvailableSlots, "id = ?", input.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No existing record found for the given ID"})
		return
	}

	// Update only the fields provided in the payload
	if input.DayOfWeek != "" {
		existingAvailableSlots.DayOfWeek = input.DayOfWeek
	}
	if !input.StartTime.IsZero() {
		existingAvailableSlots.StartTime = input.StartTime
	}
	if !input.EndTime.IsZero() {
		existingAvailableSlots.EndTime = input.EndTime
	}
	existingAvailableSlots.IsRecurring = input.IsRecurring // Update bool field, default `false` when not sent
	if input.CustomDate != nil {
		existingAvailableSlots.CustomDate = input.CustomDate
	}

	// Save the updated record back to the database
	if err := config.DB.Save(&existingAvailableSlots).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update available slots"})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Available slots updated successfully",
		"data":    existingAvailableSlots,
	})
}

func GetAvailableSlots(c *gin.Context) {
	var availabileSlots []models.AvailableTimeSlots

	// Query parameters for optional filtering
	UserID := c.Query("user_id")
	dayOfweek := c.Query("day_of_week")

	query := config.DB
	if UserID != "" {
		query = query.Where("user_id", UserID)
	}

	if dayOfweek != "" {
		query = query.Where("day_of_week=?", dayOfweek)
	}
	if err := query.Find(&availabileSlots).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve available slots"})
		return
	}
	//Check if recors exist
	if len(availabileSlots) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No available slots found"})
	}
	//Success response

	c.JSON(http.StatusOK, gin.H{
		"message": "Available slotrs rereived sucessfully",
		"data":    availabileSlots,
	})
}

func DeleteAvailableSlot(c *gin.Context) {
	var slot models.AvailableTimeSlots

	// Extract user_id and slot_id from query parameters
	userID := c.Param("user_id")
	slotID := c.Param("slot_id")

	// Validate parameters
	if userID == "" || slotID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID and Slot ID are required"})
		return
	}

	// Find the slot by ID and UserID
	if err := config.DB.Where("id = ? AND user_id = ?", slotID, userID).First(&slot).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Slot not found or does not belong to the user"})
		return
	}

	// Delete the slot
	if err := config.DB.Delete(&slot).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the slot"})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Slot deleted successfully",
		"slot_id": slotID,
		"user_id": userID,
	})
}
