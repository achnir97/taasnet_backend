package handlers

import (
	"log"
	"net/http"
	"strings"
	"taas-api/config"
	"taas-api/models"
	"taas-api/services"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AvailabilityRequest struct {
	TalentID       string   `json:"talent_id" binding:"required"`
	AvailableDate  string   `json:"available_date" binding:"required"` // Format: YYYY-MM-DD
	AvailableSlots []string `json:"available_slots" binding:"required"`
}

func removeDuplicates(slots []string) []string {
	slotMap := make(map[string]bool)
	var uniqueSlots []string
	for _, slot := range slots {
		if !slotMap[slot] {
			slotMap[slot] = true
			uniqueSlots = append(uniqueSlots, slot)
		}
	}
	return uniqueSlots
}

func CreateAvailableSlots(c *gin.Context) {
	var req AvailabilityRequest

	// Parse and validate the request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Parse the date
	availableDate, err := time.Parse("2006-01-02", req.AvailableDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD."})
		return
	}

	// Check if an entry for this talent and date already exists
	var existingRecord models.AvailableTimeSlots
	err = config.DB.Where("talent_id = ? AND available_date = ?", req.TalentID, availableDate).
		First(&existingRecord).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database."})
		return
	}

	if err == nil {
		// Record exists, update the slots
		existingSlots := strings.Split(existingRecord.AvailableSlots, ",")
		updatedSlots := append(existingSlots, req.AvailableSlots...)
		uniqueSlots := removeDuplicates(updatedSlots)
		existingRecord.AvailableSlots = strings.Join(uniqueSlots, ",")

		// Use `Updates` to ensure proper WHERE conditions
		if err := config.DB.Model(&models.AvailableTimeSlots{}).
			Where("talent_id = ? AND available_date = ?", req.TalentID, availableDate).
			Updates(map[string]interface{}{
				"available_slots": existingRecord.AvailableSlots,
				"updated_at":      time.Now(),
			}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update availability."})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Availability updated successfully",
			"data":    existingRecord,
		})
		return
	}

	// Record doesn't exist, create a new entry
	newRecord := models.AvailableTimeSlots{
		TalentID:       req.TalentID,
		AvailableDate:  availableDate,
		AvailableSlots: strings.Join(req.AvailableSlots, ","),
	}

	if err := config.DB.Create(&newRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create availability record."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Availability created successfully",
		"data":    newRecord,
	})
}

func EditAvailableSlot(c *gin.Context) {
	var req struct {
		TalentID string `json:"talent_id" binding:"required"`
		Date     string `json:"date" binding:"required"`     // Format: YYYY-MM-DD
		OldSlot  string `json:"old_slot" binding:"required"` // Format: HH:mm-HH:mm
		NewSlot  string `json:"new_slot" binding:"required"` // Format: HH:mm-HH:mm
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Parse the date
	availableDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD."})
		return
	}

	// Fetch the record for the given date and talent
	var existingRecord models.AvailableTimeSlots
	err = config.DB.Where("talent_id = ? AND available_date = ?", req.TalentID, availableDate).
		First(&existingRecord).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Availability record not found."})
		return
	}

	// Update the slot
	slots := strings.Split(existingRecord.AvailableSlots, ",")
	for i, slot := range slots {
		if slot == req.OldSlot {
			slots[i] = req.NewSlot
			existingRecord.AvailableSlots = strings.Join(slots, ",")
			if err := config.DB.Save(&existingRecord).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update slot."})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Slot updated successfully", "data": existingRecord})
			return
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Old slot not found."})
}

func DeleteAvailableSlot(c *gin.Context) {
	var req struct {
		TalentID string `json:"talent_id" binding:"required"`
		Date     string `json:"date" binding:"required"` // Format: YYYY-MM-DD
		Slot     string `json:"slot" binding:"required"` // Format: HH:mm-HH:mm
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Parse the date
	availableDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD."})
		return
	}

	// Fetch the record for the given date and talent
	var existingRecord models.AvailableTimeSlots
	err = config.DB.Where("talent_id = ? AND available_date = ?", req.TalentID, availableDate).
		First(&existingRecord).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Availability record not found."})
		return
	}

	// Remove the slot
	slots := strings.Split(existingRecord.AvailableSlots, ",")
	newSlots := []string{}
	for _, slot := range slots {
		if slot != req.Slot {
			newSlots = append(newSlots, slot)
		}
	}

	if len(newSlots) == len(slots) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slot not found."})
		return
	}

	existingRecord.AvailableSlots = strings.Join(newSlots, ",")
	if err := config.DB.Save(&existingRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete slot."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Slot deleted successfully", "data": existingRecord})
}

func GetTalentAvailability(c *gin.Context) {
	talentID := c.Query("talent_id")
	if talentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Talent ID is required"})
		return
	}

	// Query the database for available time slots
	var availability []models.AvailableTimeSlots
	if err := config.DB.Where("talent_id = ?", talentID).Find(&availability).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch availability"})
		return
	}

	// Define the cutoff time (1 hour ahead of the current time)
	cutoffTime := time.Now().Add(1 * time.Hour)
	log.Printf("Cutoff time: %v", cutoffTime)

	// Filter and organize data
	filteredAvailability := services.FilterValidSlots(availability, cutoffTime)

	// Format the response
	availabilityMap := make([]map[string]interface{}, 0)

	for _, record := range filteredAvailability {
		availabilityMap = append(availabilityMap, map[string]interface{}{
			"date":            record.AvailableDate.Format("2006-01-02"),
			"available_slots": strings.Split(record.AvailableSlots, ","),
		})
	}

	// Send optimized response
	c.JSON(http.StatusOK, gin.H{
		"talent_id":    talentID,
		"availability": availabilityMap,
	})
}
