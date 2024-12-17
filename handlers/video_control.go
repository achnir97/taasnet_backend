package handlers

import (
	"fmt"
	"net/http"
	"taas-api/config"
	"taas-api/models"

	"github.com/gin-gonic/gin"
)

func SaveVideoControl(c *gin.Context) {
	var input models.VideoControl
	var existingVideoControl models.VideoControl

	// Parse the input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if a record already exists
	if err := config.DB.First(&existingVideoControl).Error; err == nil {
		// Update the existing record with the new data
		existingVideoControl.VideoURL = input.VideoURL
		existingVideoControl.Action = input.Action
		existingVideoControl.StartTime = input.StartTime
		existingVideoControl.EndTime = input.EndTime
		existingVideoControl.PausedTime = input.PausedTime

		// Save the updated record
		if err := config.DB.Save(&existingVideoControl).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update video control"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":       "Video control updated successfully",
			"video_control": existingVideoControl,
		})
		return
	}

	// If no record exists, create a new one
	newVideoControl := models.VideoControl{
		VideoURL:   input.VideoURL,
		Action:     input.Action,
		StartTime:  input.StartTime,
		EndTime:    input.EndTime,
		PausedTime: input.PausedTime,
	}

	if err := config.DB.Create(&newVideoControl).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save video control"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message":       "New video control created successfully",
		"video_control": newVideoControl,
	})
}

// GetVideoControl handles fetching the single video control record
func GetVideoControl(c *gin.Context) {
	// Create an instance to store the video control data
	var videoControl models.VideoControl

	// Fetch the first video control record from the database
	if err := config.DB.First(&videoControl).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch video control"})
		return
	}

	// Return success response with video control data
	c.JSON(http.StatusOK, gin.H{
		"message":       "Video control fetched successfully",
		"video_control": videoControl,
	})
}

func UpdateVideoControl(c *gin.Context) {
	var input models.VideoControl
	var existingVideoControl models.VideoControl

	// Parse the input JSON
	// Parse the input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("Error binding JSON:", err) // Log the error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("Received JSON payload:", input) // Debugging: Log received data

	// Check if a record exists
	if err := config.DB.First(&existingVideoControl).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No existing video control record found"})
		return
	}

	// Update only the fields provided
	if input.VideoURL != "" {
		existingVideoControl.VideoURL = input.VideoURL
	}
	if input.StartTime > 0 {
		existingVideoControl.StartTime = input.StartTime
	}
	if input.EndTime > 0 {
		existingVideoControl.EndTime = input.EndTime
	}
	if input.PausedTime >= 0 {
		existingVideoControl.PausedTime = input.PausedTime
	}
	if input.Action != "" {
		existingVideoControl.Action = input.Action
	}

	// Save the updated record
	if err := config.DB.Save(&existingVideoControl).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update video control"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Video control updated successfully",
		"video_control": existingVideoControl,
	})
}
