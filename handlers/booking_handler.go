package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"taas-api/config"
	"taas-api/models"
	"time"

	"github.com/gin-gonic/gin"
	"storj.io/common/uuid"
)

type CreateBookingRequest struct {
	CardID          string `json:"card_id" binding:"required"`        // Card ID
	CardTitle       string `json:"card_title" binding:"required"`     // Card Title
	UserID          string `json:"user_id" binding:"required"`        // User ID
	TalentID        string `json:"talent_id" binding:"required"`      // Talent ID
	SessionType     string `json:"session_type" binding:"required"`   // Enum: CoffeeCall, Regular
	Status          string `json:"status" binding:"required"`         // Enum: Scheduled, Completed, Cancelled
	PaymentStatus   string `json:"payment_status" binding:"required"` // Enum: Paid, Pending
	SpecialRequests string `json:"special_requests,omitempty"`        // Optional
	CardDuration    int    `json:"card_duration" binding:"required"`  // Card Duration (in minutes)

	Slots []struct {
		BookingDate string   `json:"booking_date" binding:"required"`
		TimeSlots   []string `json:"time_slots" binding:"required"`
	} `json:"slots" binding:"required"`
}

// Handler for POST /bookings to create a new booking
func CreateBookings(c *gin.Context) {
	var req CreateBookingRequest

	// Clone the request body for logging
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return
	}
	// Log the raw request body
	fmt.Println("Request Body:", string(body))

	// Reset the request body for binding
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// Bind JSON data from the request body to the request struct
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Parsed Request Data:", req)

	// Iterate through each slot and create a booking record
	for _, slot := range req.Slots {
		bookingDate, err := time.Parse("2006-01-02", slot.BookingDate)
		if err != nil {
			fmt.Println("Invalid booking_date format:", slot.BookingDate, "Error:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking_date format"})
			return
		}

		var timeRanges models.TimeRanges
		//iterate through all the time slots for the date
		for _, timeSlot := range slot.TimeSlots {
			startTime, endTime, err := parseTimeRange(timeSlot)
			if err != nil {
				fmt.Println("Invalid time_slot format:", timeSlot, "Error:", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time_slot format"})
				return
			}
			parsedStartTime, err := time.Parse("15:04", startTime)
			if err != nil {
				fmt.Println("Error parsing start time", startTime, err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format"})
				return
			}
			parsedEndTime, err := time.Parse("15:04", endTime)
			if err != nil {
				fmt.Println("Error parsing end time", endTime, err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time format"})
				return
			}
			timeRanges = append(timeRanges, models.TimeRange{
				StartTime: parsedStartTime.Format("15:04"),
				EndTime:   parsedEndTime.Format("15:04"),
			})
		}

		// Generate a new UUID for the BookingID
		bookingID, err := uuid.New()
		if err != nil {
			fmt.Println("Error generating BookingID:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate booking ID"})
			return
		}
		fmt.Println("Generated BookingID:", bookingID.String())

		// Map the request data to the Bookings model
		newBooking := models.BookingRequests{
			BookingID:       bookingID.String(),
			CardID:          req.CardID,
			CardTitle:       req.CardTitle,
			UserID:          req.UserID,
			TalentID:        req.TalentID,
			SessionType:     models.SessionType(req.SessionType),
			BookedTime:      timeRanges,
			Status:          models.BookingStatus(req.Status),
			PaymentStatus:   models.PaymentStatus(req.PaymentStatus),
			SpecialRequests: req.SpecialRequests,
			BookingDate:     bookingDate,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		fmt.Println("New Booking Data:", newBooking)
		// Create the booking in the database
		if err := config.DB.Create(&newBooking).Error; err != nil {
			fmt.Println("Database Error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
			return
		}

		fmt.Println("Booking created successfully:", newBooking)

	}
	// Respond with a success message
	c.JSON(http.StatusCreated, gin.H{"message": "Bookings created successfully"})
}
func parseTimeRange(timeRange string) (string, string, error) {
	parts := strings.Split(timeRange, "-")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid time range format")
	}
	return parts[0], parts[1], nil
}

// Handler for GET /bookings/user/:user_id to fetch all bookings made by a user
func GetBookingsByUser(c *gin.Context) {
	log.Println("Starting GetBookingsByTalent handler")

	// Extract talent_id from URL params
	userId := c.Param("user_id")
	log.Printf("Extracted talent_id from URL: %s\n", userId)

	if userId == "" {
		log.Println("Error: Talent ID is required in the URL")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Talent ID is required"})
		return
	}

	// Fetch bookings from the database based on talent_id
	log.Printf("Fetching bookings from database for talent_id: %s\n", userId)
	var bookings []models.BookingRequests
	result := config.DB.Where("user_id = ?", userId).Find(&bookings)

	if result.Error != nil {
		log.Printf("Database Error for talent_id: %s, Error: %v\n", userId, result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}
	log.Printf("Successfully fetched %d bookings from database for user_id: %s\n", len(bookings), userId)

	// Log the fetched bookings and respond
	log.Printf("Responding with bookings: %v for user_id: %s\n", bookings, userId)
	c.JSON(http.StatusOK, gin.H{"bookings": bookings})
	log.Println("Successfully completed GetBookingsByTalent handler")
}

// Handler for GET /bookings/talent/:talent_id to fetch all booking requests for a talent
func GetBookingsByTalent(c *gin.Context) {
	log.Println("Starting GetBookingsByTalent handler")

	// Extract talent_id from URL params
	talentID := c.Param("talent_id")
	log.Printf("Extracted talent_id from URL: %s\n", talentID)

	if talentID == "" {
		log.Println("Error: Talent ID is required in the URL")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Talent ID is required"})
		return
	}

	// Fetch bookings from the database based on talent_id
	log.Printf("Fetching bookings from database for talent_id: %s\n", talentID)
	var bookings []models.BookingRequests
	result := config.DB.Where("talent_id = ?", talentID).Find(&bookings)

	if result.Error != nil {
		log.Printf("Database Error for talent_id: %s, Error: %v\n", talentID, result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}
	log.Printf("Successfully fetched %d bookings from database for talent_id: %s\n", len(bookings), talentID)

	// Log the fetched bookings and respond
	log.Printf("Responding with bookings: %v for talent_id: %s\n", bookings, talentID)
	c.JSON(http.StatusOK, gin.H{"bookings": bookings})
	log.Println("Successfully completed GetBookingsByTalent handler")
}
