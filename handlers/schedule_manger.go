package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"taas-api/config"
	"taas-api/models"

	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BookingRequest struct to handle data from frontend
type BookingRequest struct {
	TalentID  string    `json:"talent_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
	UserID    string    `json:"user_id" binding:"required"`
}

type AvailabilityRequest struct {
	TalentID       string            `json:"talent_id" binding:"required"`
	AvailableDate  string            `json:"available_date" binding:"required"`
	AvailableSlots models.TimeRanges `json:"available_slots" binding:"required"` // Using TimeRanges type
}

// Handler to create or update available time slots
func CreateAvailableSlots(c *gin.Context) {
	log.Println("Starting CreateAvailableSlots handler")

	var req AvailabilityRequest

	// Parse and validate the request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error parsing input JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}
	log.Printf("Parsed request: %+v\n", req)

	// Parse the date
	availableDate, err := time.Parse("2006-01-02", req.AvailableDate)
	if err != nil {
		log.Printf("Error parsing date: %s, Error: %v\n", req.AvailableDate, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD."})
		return
	}
	log.Printf("Parsed date: %s\n", availableDate.Format("2006-01-02"))

	// Check if an entry for this talent and date already exists
	var existingRecord models.AvailableTimeSlots
	log.Printf("Querying database for talent_id: %s, date: %s\n", req.TalentID, availableDate)
	err = config.DB.Where("talent_id = ? AND available_date = ?", req.TalentID, availableDate).
		First(&existingRecord).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		log.Printf("Database query error for talent_id: %s, date: %s, Error: %v\n", req.TalentID, availableDate, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database."})
		return
	}
	if err == nil {
		// Record exists, update the slots
		log.Printf("Existing record found: %+v\n", existingRecord)
		//Check for booking collisions.
		log.Println("Checking for booking collisions")
		if !areSlotsAvailable(config.DB, req.TalentID, availableDate, req.AvailableSlots) {
			log.Println("Booking collision detected, cannot save.")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Some of the available time slots you are trying to add are already booked, try again with different time."})
			return
		}

		updatedSlots := append(existingRecord.AvailableSlots, req.AvailableSlots...)
		log.Printf("Updated slots: %v\n", updatedSlots)

		existingRecord.AvailableSlots = updatedSlots
		log.Println("Updating availability record in database")
		if err := config.DB.Save(&existingRecord).Error; err != nil {
			log.Printf("Error updating availability: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update availability."})
			return
		}
		log.Printf("Availability updated successfully: %+v\n", existingRecord)
		c.JSON(http.StatusOK, gin.H{
			"message": "Availability updated successfully",
			"data":    existingRecord,
		})
		log.Println("Successfully updated and returning result")
		return
	}
	// Record doesn't exist, create a new entry
	log.Println("No existing record found")
	if !areSlotsAvailable(config.DB, req.TalentID, availableDate, req.AvailableSlots) {
		log.Println("Booking collision detected, cannot save.")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some of the available time slots you are trying to add are already booked, try again with different time."})
		return
	}
	newRecord := models.AvailableTimeSlots{
		TalentID:       req.TalentID,
		AvailableDate:  availableDate,
		AvailableSlots: req.AvailableSlots,
	}
	log.Printf("Creating new availability record: %+v\n", newRecord)
	if err := config.DB.Create(&newRecord).Error; err != nil {
		log.Printf("Error creating availability record: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create availability record."})
		return
	}
	log.Printf("Availability created successfully: %+v\n", newRecord)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Availability created successfully",
		"data":    newRecord,
	})
	log.Println("Successfully created and returning result")
}

func FetchrawAllAvailableTimeSlots(c *gin.Context) {
	// Assuming you pass DB as middleware.
	talentID := c.Query("talent_id")

	if talentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Talent ID is required"})
		return
	}

	var availableTimeSlots []models.AvailableTimeSlots
	log.Println("Fetching availabilities for talent id : ", talentID)
	if err := config.DB.Where("talent_id = ?", talentID).Find(&availableTimeSlots).Error; err != nil {
		log.Printf("Database error while fetching talent availability: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch talent availabilities: " + err.Error()})
		return
	}
	log.Printf("Successfully fetched availabilities for talent_id: %s,  records count: %d\n", talentID, len(availableTimeSlots))
	response := gin.H{"available_slots": availableTimeSlots}

	responseBytes, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		log.Printf("Error while marshaling response: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to format response data"})
		return
	}
	log.Println("Response data:")
	log.Println(string(responseBytes))

	c.JSON(http.StatusOK, response)

}
