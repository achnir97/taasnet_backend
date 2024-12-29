package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"taas-api/config"
	"taas-api/models"

	"io"
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
	StartTime       string `json:"start_time" binding:"required"`     // ISO8601 format, parsed in handler
	EndTime         string `json:"end_time" binding:"required"`       // ISO8601 format, parsed in handler
	Status          string `json:"status" binding:"required"`         // Enum: Scheduled, Completed, Cancelled
	PaymentStatus   string `json:"payment_status" binding:"required"` // Enum: Paid, Pending
	SpecialRequests string `json:"special_requests,omitempty"`        // Optional
}

// NormalizeTime ensures the time string is in the proper HH:mm format
func NormalizeTime(timeStr string) string {
	parts := strings.Split(timeStr, ":")
	if len(parts) == 2 && len(parts[1]) == 1 { // If minutes are a single digit
		return fmt.Sprintf("%s:0%s", parts[0], parts[1])
	}
	return timeStr
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

	// Normalize the start_time and end_time
	req.StartTime = NormalizeTime(req.StartTime)
	req.EndTime = NormalizeTime(req.EndTime)

	// Parse StartTime and EndTime to time.Time
	startTime, err := time.Parse("2006-01-02 15:04", req.StartTime)
	if err != nil {
		fmt.Println("Invalid start_time format:", req.StartTime, "Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format"})
		return
	}
	fmt.Println("Parsed StartTime:", startTime)

	endTime, err := time.Parse("2006-01-02 15:04", req.EndTime)
	if err != nil {
		fmt.Println("Invalid end_time format:", req.EndTime, "Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format"})
		return
	}
	fmt.Println("Parsed EndTime:", endTime)

	// Generate a new UUID for the BookingID
	bookingID, err := uuid.New()
	if err != nil {
		fmt.Println("Error generating BookingID:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate booking ID"})
		return
	}
	fmt.Println("Generated BookingID:", bookingID.String())

	// Map the request data to the Bookings model
	newBooking := models.BookingRequest{
		BookingID:       bookingID.String(), // Generate booking ID in the backend
		CardID:          req.CardID,
		CardTitle:       req.CardTitle,
		UserID:          req.UserID,
		TalentID:        req.TalentID,
		SessionType:     req.SessionType,
		StartTime:       startTime,
		EndTime:         endTime,
		Status:          req.Status,
		PaymentStatus:   req.PaymentStatus,
		SpecialRequests: req.SpecialRequests,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Log the final booking data before saving to the database
	fmt.Println("New Booking Data:", newBooking)

	// Create the booking in the database
	result := config.DB.Create(&newBooking)
	if result.Error != nil {
		fmt.Println("Database Error:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	// Log success and respond
	fmt.Println("Booking created successfully:", newBooking)
	c.JSON(http.StatusCreated, newBooking)
}

// Handler for GET /bookings/user/:user_id to fetch all bookings made by a user
func GetBookingsByUser(c *gin.Context) {
	// Extract user_id from URL params
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Fetch bookings from the database based on user_id
	var bookings []models.BookingRequest
	result := config.DB.Where("user_id = ?", userID).Find(&bookings)

	if result.Error != nil {
		fmt.Println("Database Error:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	// Log the fetched bookings and respond
	fmt.Println("Fetched Bookings for User ID:", userID, bookings)
	c.JSON(http.StatusOK, bookings)
}

// Handler for GET /bookings/talent/:talent_id to fetch all booking requests for a talent
func GetBookingsByTalent(c *gin.Context) {
	// Extract talent_id from URL params
	talentID := c.Param("talent_id")
	if talentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Talent ID is required"})
		return
	}

	// Fetch bookings from the database based on talent_id
	var bookings []models.BookingRequest
	result := config.DB.Where("talent_id = ?", talentID).Find(&bookings)

	if result.Error != nil {
		fmt.Println("Database Error:", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	// Log the fetched bookings and respond
	fmt.Println("Fetched Bookings for Talent ID:", talentID, bookings)
	c.JSON(http.StatusOK, bookings)
}
