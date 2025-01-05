package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"taas-api/config"
	"taas-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Response struct for frontend
type AvailableSlotResponse struct {
	AvailableSlots map[string][]string `json:"available_slots"`
	BookedSlots    map[string][]string `json:"booked_slots"`
}

// Handler to fetch all the available slots for a specific date.
func FetchBookFilteredAvailableTimeSlots(c *gin.Context) {
	db := config.DB
	talentID := c.Query("talent_id")
	cardDurationParam := c.Query("duration")

	log.Printf("Starting FetchFilteredAvailableTimeSlots with talentID: %s, duration: %s\n", talentID, cardDurationParam)

	if talentID == "" || cardDurationParam == "" {
		log.Println("Error: Talent ID and Duration are required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Talent ID and Duration is required"})
		return
	}

	duration, err := time.ParseDuration(cardDurationParam + "m")
	if err != nil {
		log.Printf("Error: Invalid duration format: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid duration format"})
		return
	}
	log.Printf("Parsed duration: %v\n", duration)
	// Calculate the start of today (midnight)
	today := time.Now().Truncate(time.Hour * 24)

	var availableTimeSlots []models.AvailableTimeSlots
	log.Println("Fetching all available time slots")
	// Fetch available slots that are today or in future
	if err := db.Where("talent_id = ? AND available_date >= ?", talentID, today).Find(&availableTimeSlots).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("Error: Talent is not availabe on this selected date. ")
			c.JSON(http.StatusNotFound, gin.H{"error": "Talent is not availabe on this selected date. "})
			return
		}
		log.Printf("Error: Failed to fetch talent availability: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch talent availability. " + err.Error()})
		return
	}

	log.Printf("Successfully fetched available time slots. Number of records: %d \n", len(availableTimeSlots))

	availableSlotsMap := make(map[string][]string)
	bookedSlotsMap := make(map[string][]string)

	for _, availableTimeSlot := range availableTimeSlots {
		log.Printf("Processing available time slot for date: %v\n", availableTimeSlot.AvailableDate)

		// Fetch booked slots for the date
		var bookings []models.BookingRequests
		log.Printf("Fetching bookings for talent id: %s and date: %v\n", talentID, availableTimeSlot.AvailableDate)
		if err := db.Where("talent_id = ? AND booking_date = ?", talentID, availableTimeSlot.AvailableDate).Find(&bookings).Error; err != nil {
			log.Printf("Error: Failed to fetch booked slots: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch booked slots: " + err.Error()})
			return
		}
		log.Printf("Successfully fetched booked slots. Number of records: %d\n", len(bookings))
		// Print all fetched booked slots
		for _, booking := range bookings {
			log.Printf("   Fetched Booked Slots: %+v\n", booking)
		}

		// Generate all possible slots based on the first available time range. If there are multiple ranges, combine them
		var allSlots []string
		for _, timeRange := range availableTimeSlot.AvailableSlots {
			startTime, _ := time.Parse("15:04", timeRange.StartTime)
			endTime, _ := time.Parse("15:04", timeRange.EndTime)
			allSlots = append(allSlots, generateAllTimeSlots(startTime, endTime, duration)...)

		}
		log.Printf("  Generated all possible slots: %v\n", allSlots)

		// Collect all booked time ranges
		bookedTimeRanges := make([]models.TimeRange, 0)
		for _, booking := range bookings {
			bookedTimeRanges = append(bookedTimeRanges, booking.BookedTime...)
		}
		log.Printf("   Booked time ranges from booking table: %v\n", bookedTimeRanges)

		var availableSlots []string
		var bookedSlots []string

		// Compare and generate available and booked slots
		if availableTimeSlot.AvailableDate.Equal(time.Now().Truncate(time.Hour * 24)) {
			currentTime := time.Now()
			for _, slot := range allSlots {
				slotStartTime, slotEndTime, err := parseTimeRange1(slot)
				if err != nil {
					log.Printf("Error parsing time range: %v\n", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse time range: " + err.Error()})
					return
				}
				parsedSlotStartTime, err := time.Parse("15:04", slotStartTime)
				if err != nil {
					log.Printf("Error parsing time: %v\n", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse time: " + err.Error()})
					return
				}

				slotTime := time.Date(availableTimeSlot.AvailableDate.Year(), availableTimeSlot.AvailableDate.Month(), availableTimeSlot.AvailableDate.Day(), parsedSlotStartTime.Hour(), parsedSlotStartTime.Minute(), 0, 0, time.Local)
				if slotTime.After(currentTime) {
					if isTimeSlotOverlapping(slotStartTime, slotEndTime, bookedTimeRanges) {
						bookedSlots = append(bookedSlots, slot)
					} else {
						availableSlots = append(availableSlots, slot)
					}

				} else {
					if isTimeSlotOverlapping(slotStartTime, slotEndTime, bookedTimeRanges) {
						bookedSlots = append(bookedSlots, slot)
					}
					// Slots in the past are skipped, they are not available now.
				}

			}
		} else {
			for _, slot := range allSlots {
				slotStartTime, slotEndTime, err := parseTimeRange(slot)
				if err != nil {
					log.Printf("Error parsing time range: %v\n", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse time range: " + err.Error()})
					return
				}
				if isTimeSlotOverlapping(slotStartTime, slotEndTime, bookedTimeRanges) {
					bookedSlots = append(bookedSlots, slot)
				} else {
					availableSlots = append(availableSlots, slot)
				}

			}
		}
		log.Printf("  Available slots after compare: %v\n", availableSlots)
		log.Printf("  Booked slots after compare: %v\n", bookedSlots)

		// Store slots in map for each date
		availableSlotsMap[availableTimeSlot.AvailableDate.Format("2006-01-02")] = availableSlots
		bookedSlotsMap[availableTimeSlot.AvailableDate.Format("2006-01-02")] = bookedSlots

		log.Printf("Completed processing for available time slot for date: %v\n", availableTimeSlot.AvailableDate)

	}

	// Create response struct
	response := AvailableSlotResponse{
		AvailableSlots: availableSlotsMap,
		BookedSlots:    bookedSlotsMap,
	}

	c.JSON(http.StatusOK, response)
	log.Println("Successfully returned all available slots.")

}

// Helper function to generate all possible slots based on start time, end time and card duration
func generateAllTimeSlots(startTime time.Time, endTime time.Time, duration time.Duration) []string {
	var slots []string
	current := startTime
	for current.Before(endTime) {
		next := current.Add(duration)
		slots = append(slots, current.Format("15:04")+"-"+next.Format("15:04"))
		current = next
	}

	return slots
}
func parseTimeRange1(timeRange string) (string, string, error) {
	parts := strings.Split(timeRange, "-")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid time range format")
	}
	return parts[0], parts[1], nil
}

func isTimeSlotOverlapping(slotStartTime string, slotEndTime string, bookedTimeRanges []models.TimeRange) bool {

	parsedSlotStartTime, _ := time.Parse("15:04", slotStartTime)
	parsedSlotEndTime, _ := time.Parse("15:04", slotEndTime)

	// Iterate through all booked time ranges to check for overlap
	for _, bookedTimeRange := range bookedTimeRanges {
		parsedBookedStartTime, _ := time.Parse("15:04", bookedTimeRange.StartTime)
		parsedBookedEndTime, _ := time.Parse("15:04", bookedTimeRange.EndTime)

		// Check if the current slot overlaps with the booked time range
		if !(parsedSlotEndTime.Before(parsedBookedStartTime) || parsedSlotStartTime.After(parsedBookedEndTime)) {
			return true // Overlapping
		}
	}
	return false
}
