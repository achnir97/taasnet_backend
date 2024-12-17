package handlers

import (
	"net/http"
	"taas-api/config"
	_ "taas-api/config"
	"taas-api/models"

	"github.com/gin-gonic/gin"
)

func CreateTimeSlot(c *gin.Context) {
	var slot models.AvailableTimeSlots

	// Bind the request body to the struct
	if err := c.ShouldBindJSON(&slot); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save to the database
	if err := config.DB.Create(&slot).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create time slot"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Time slot created successfully", "slot": slot})
}

func GetUserTimeSlots(c *gin.Context) {
	var slots []models.AvailableTimeSlots
	userID := c.Param("user_id")

	// Fetch time slots by UserID
	if err := config.DB.Where("user_id = ?", userID).Find(&slots).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch time slots"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"time_slots": slots})
}

func DeleteTimeSlot(c *gin.Context) {
	slotID := c.Param("id")

	// Soft delete the time slot
	if err := config.DB.Delete(&models.AvailableTimeSlots{}, slotID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete time slot"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Time slot deleted successfully"})
}
