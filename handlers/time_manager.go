package handlers

import (
	"fmt"
	"log"
	"net/http"
	"taas-api/config"
	"taas-api/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler to fetch all the available slots for a specific date.
func FetchFilteredAvailableTimeSlots(c *gin.Context) {
	db := config.DB
	talentID := c.Query("talent_id")
	cardDuration := c.Query("duration")
	log.Printf("Starting FetchFilteredAvailableTimeSlots with talentID: %s, duration: %s\n", talentID, cardDuration)

	if talentID == "" || cardDuration == "" {
		log.Println("Error: Talent ID and Duration are required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Talent ID and Duration is required"})
		return
	}

	var availableTimeSlots []models.AvailableTimeSlots
	log.Println("Fetching all available time slots")
	if err := db.Where("talent_id = ? ", talentID).Find(&availableTimeSlots).Error; err != nil {
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

	duration, err := time.ParseDuration(cardDuration + "m")
	if err != nil {
		log.Printf("Error: Invalid duration format: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid duration format"})
		return
	}
	log.Printf("Parsed duration: %v\n", duration)

	allAvailableSlots := make(map[string][]map[string]string)
	log.Println("Starting iteration through all available time slots")

	for _, availableTimeSlot := range availableTimeSlots {
		log.Printf("Processing available time slot for date: %v\n", availableTimeSlot.AvailableDate)

		// Check if the date is in the future!
		var bookings []models.BookingRequests
		log.Printf("Fetching bookings for talent id: %s and date: %v\n", talentID, availableTimeSlot.AvailableDate)
		if err := db.Where("talent_id = ? AND booking_date = ?", talentID, availableTimeSlot.AvailableDate).Find(&bookings).Error; err != nil {
			log.Printf("Error: Failed to fetch booked slots: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch booked slots: " + err.Error()})
			return
		}
		log.Printf("Successfully fetched booked slots. Number of records: %d\n", len(bookings))
		for i, booking := range bookings {
			fmt.Printf("Booking [%d]: %+v\n", i, booking)
			log.Printf("Booking [%d]: %+v\n", i, booking)
		}

		if availableTimeSlot.AvailableDate.After(time.Now()) {
			log.Println("Available date is in the future. Generating available slots.")
			availableSlots := generateAvailableSlots(availableTimeSlot.AvailableSlots, bookings, duration, availableTimeSlot.AvailableDate)
			log.Printf("Successfully generated available slots: %v\n", availableSlots)
			for date, slots := range availableSlots {
				allAvailableSlots[date] = slots
			}

		} else if availableTimeSlot.AvailableDate.Equal(time.Now().Truncate(time.Hour * 24)) {
			log.Println("Available date is today. Generating available slots with cutoff.")
			currentTime := time.Now()
			availableSlots := generateAvailableSlotsWithCutoff(availableTimeSlot.AvailableSlots, bookings, duration, currentTime, availableTimeSlot.AvailableDate)
			log.Printf("Successfully generated available slots with cutoff: %v\n", availableSlots)
			for date, slots := range availableSlots {
				allAvailableSlots[date] = slots
			}
		} else {
			log.Println("Available date is in the past, skipping.")
		}
		log.Printf("Completed processing for available time slot for date: %v\n", availableTimeSlot.AvailableDate)

	}
	log.Println("Combining available slots into a final list")
	var finalAvailableSlots []map[string]interface{}
	for date, slots := range allAvailableSlots {
		finalAvailableSlots = append(finalAvailableSlots, map[string]interface{}{"date": date, "slots": slots})
	}

	log.Printf("Combined all available slots. Number of unique dates: %d\n", len(finalAvailableSlots))
	c.JSON(http.StatusOK, gin.H{"available_slots": finalAvailableSlots})
	log.Println("Successfully returned all available slots.")
}

func generateAvailableSlotsWithCutoff(timeRanges []models.TimeRange, bookings []models.BookingRequests, cardDuration time.Duration, cutoffTime time.Time, availableDate time.Time) map[string][]map[string]string {
	slotsByDate := make(map[string][]map[string]string)
	formattedDate := availableDate.Format("2006-01-02")
	var slots []map[string]string
	for _, timeRange := range timeRanges {
		startTime, _ := time.Parse("15:04", timeRange.StartTime)
		endTime, _ := time.Parse("15:04", timeRange.EndTime)
		for currentTime := startTime; currentTime.Before(endTime); currentTime = currentTime.Add(cardDuration) {
			slotEndTime := currentTime.Add(cardDuration)
			if slotEndTime.After(endTime) {
				break
			}
			slotStartTime := time.Date(availableDate.Year(), availableDate.Month(), availableDate.Day(), currentTime.Hour(), currentTime.Minute(), 0, 0, availableDate.Location())
			if !slotStartTime.Before(cutoffTime) {
				slotStart := time.Date(availableDate.Year(), availableDate.Month(), availableDate.Day(), currentTime.Hour(), currentTime.Minute(), 0, 0, availableDate.Location()).Format("15:04")
				slotEnd := time.Date(availableDate.Year(), availableDate.Month(), availableDate.Day(), slotEndTime.Hour(), slotEndTime.Minute(), 0, 0, availableDate.Location()).Format("15:04")
				// Check against booked slots
				if !isSlotBooked(bookings, slotStartTime, slotEndTime) {
					slots = append(slots, map[string]string{
						"startTime": slotStart,
						"endTime":   slotEnd,
					})
				}
			}
		}
	}
	slotsByDate[formattedDate] = slots
	return slotsByDate
}

func generateAvailableSlots(timeRanges []models.TimeRange, bookings []models.BookingRequests, cardDuration time.Duration, availableDate time.Time) map[string][]map[string]string {
	slotsByDate := make(map[string][]map[string]string)
	formattedDate := availableDate.Format("2006-01-02")
	var slots []map[string]string

	for _, timeRange := range timeRanges {
		startTime, _ := time.Parse("15:04", timeRange.StartTime)
		endTime, _ := time.Parse("15:04", timeRange.EndTime)
		for currentTime := startTime; currentTime.Before(endTime); currentTime = currentTime.Add(cardDuration) {
			slotEndTime := currentTime.Add(cardDuration)
			if slotEndTime.After(endTime) {
				break
			}
			slotStart := time.Date(availableDate.Year(), availableDate.Month(), availableDate.Day(), currentTime.Hour(), currentTime.Minute(), 0, 0, availableDate.Location()).Format("15:04")
			slotEnd := time.Date(availableDate.Year(), availableDate.Month(), availableDate.Day(), slotEndTime.Hour(), slotEndTime.Minute(), 0, 0, availableDate.Location()).Format("15:04")
			// Check against booked slots
			if !isSlotBooked(bookings, time.Date(availableDate.Year(), availableDate.Month(), availableDate.Day(), currentTime.Hour(), currentTime.Minute(), 0, 0, availableDate.Location()),
				time.Date(availableDate.Year(), availableDate.Month(), availableDate.Day(), slotEndTime.Hour(), slotEndTime.Minute(), 0, 0, availableDate.Location())) {
				slots = append(slots, map[string]string{
					"startTime": slotStart,
					"endTime":   slotEnd,
				})
			}
		}
	}
	slotsByDate[formattedDate] = slots
	return slotsByDate
}

func isSlotBooked(bookings []models.BookingRequests, slotStartTime time.Time, slotEndTime time.Time) bool {
	for _, booking := range bookings {
		for _, bookedTime := range booking.BookedTime {
			bookingStartTime, _ := time.Parse("2006-01-02 15:04", booking.BookingDate.Format("2006-01-02")+" "+bookedTime.StartTime)
			bookingEndTime, _ := time.Parse("2006-01-02 15:04", booking.BookingDate.Format("2006-01-02")+" "+bookedTime.EndTime)

			if slotStartTime.Before(bookingEndTime) && slotEndTime.After(bookingStartTime) {
				return true
			}
		}

	}
	return false
}

func areSlotsAvailable(db *gorm.DB, talentID string, date time.Time, newSlots []models.TimeRange) bool {
	var bookings []models.BookingRequests
	if err := db.Where("talent_id = ? AND booking_date = ?", talentID, date).Find(&bookings).Error; err != nil {
		return false
	}
	for _, newSlot := range newSlots {
		newSlotStartTime, _ := time.Parse("15:04", newSlot.StartTime)
		newSlotEndTime, _ := time.Parse("15:04", newSlot.EndTime)
		for _, booking := range bookings {
			for _, bookedTime := range booking.BookedTime {
				bookingStartTime, _ := time.Parse("2006-01-02 15:04", booking.BookingDate.Format("2006-01-02")+" "+bookedTime.StartTime)
				bookingEndTime, _ := time.Parse("2006-01-02 15:04", booking.BookingDate.Format("2006-01-02")+" "+bookedTime.EndTime)
				if newSlotStartTime.Before(bookingEndTime) && newSlotEndTime.After(bookingStartTime) {
					return false
				}
			}
		}
	}
	return true
}
func isSlotAvailable(availableSlots []models.TimeRange, startTime time.Time, endTime time.Time) bool {
	start := startTime.Format("15:04")
	end := endTime.Format("15:04")
	for _, slot := range availableSlots {
		if slot.StartTime == start && slot.EndTime == end {
			return true
		}
	}
	return false
}
func removeBookedSlot(availableSlots []models.TimeRange, bookingStartTime time.Time, bookingEndTime time.Time) []models.TimeRange {
	start := bookingStartTime.Format("15:04")
	end := bookingEndTime.Format("15:04")
	for i, slot := range availableSlots {
		if slot.StartTime == start && slot.EndTime == end {
			return append(availableSlots[:i], availableSlots[i+1:]...)
		}
	}
	return nil
}
