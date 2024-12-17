package handlers

import (
	"fmt"
	"net/http"
	"taas-api/config"
	"taas-api/models"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"
	"storj.io/common/uuid"
)

// Signup Handler
func Signup(c *gin.Context) {

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Validate request body
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Check if email already exists
	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	//Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	//Generate a unique user ID using UUID
	uniqueUserID, err := uuid.New()
	if err != nil {
		fmt.Println("Error generating UUID:", err)
		return
	}

	// Create the user
	user := models.User{
		ID:       uniqueUserID.String(), // Convert to UUI
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	// Save to database
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Return success response with user ID
	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"userId":  uniqueUserID,
	})
}

// Login Handler
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Validate request body
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Fetch user by email
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	// Respond with user_id and success message
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"userId":  user.ID, // Include user ID in the response
	})
}

// / SaveCard handles the saving of card details
func SaveCard(c *gin.Context) {
	var input models.Card
	// 1. Parse and validate the incoming JSON payload
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Parse EventDate and EventTime

	// 3. Create a new Card instance
	card := models.Card{
		UserID:       input.UserID,
		Title:        input.Title,
		Description:  input.Description,
		Category:     input.Category,
		EventType:    input.EventType,
		Price:        input.Price,
		EventDate:    input.EventDate, // Already a time.Time value
		EventTime:    input.EventTime, // Already a time.Time value
		Participants: input.Participants,
		VideoURL:     input.VideoURL,
	}

	// 4. Save the card to the database
	if err := config.DB.Create(&card).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save card"})
		return
	}
	// 5. Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Card saved successfully",
		"card":    card,
	})
}

// GetUserCards handles fetching all cards for a specific user
func GetUserCards(c *gin.Context) {
	// 1. Extract UserID from query parameters
	userID := c.Query("user_id") // Assumes the user ID is sent as a query parameter

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 2. Fetch all cards for the given UserID
	var cards []models.Card
	if err := config.DB.Where("user_id = ?", userID).Find(&cards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cards"})
		return
	}

	// 3. Return the list of cards
	c.JSON(http.StatusOK, gin.H{
		"message": "Cards fetched successfully",
		"cards":   cards,
	})
}

// GetUserCards handles fetching all cards created by talents
func GetAllCards(c *gin.Context) {
	// Fetch all cards from the database without filtering by UserID
	var cards []models.Card
	if err := config.DB.Find(&cards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cards"})
		return
	}

	// Return the list of all cards
	c.JSON(http.StatusOK, gin.H{
		"message": "All cards fetched successfully",
		"cards":   cards,
	})
}

// / Cards_Id handles fetching all cards for a specific user and optional card ID
func Cards_Id(c *gin.Context) {
	// 1. Extract user_id and id from query parameters
	userID := c.Query("user_id") //Assumes the user ID is sent as a query parameter
	cardID := c.Query("id")      //Optional: Assumes the card ID is sent as an additional parameter

	// 2. Validate user_id
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// 3. Build the query dynamically
	var cards []models.Card
	query := config.DB.Where("user_id = ?", userID)

	// If cardID is provided, add it to the query
	if cardID != "" {
		query = query.Where("id = ?", cardID)
	}
	// 4. Execute the query
	if err := query.Find(&cards).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cards"})
		return
	}
	// 5. Return the result
	if len(cards) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No cards found matching the criteria"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Cards fetched successfully",
		"cards":   cards,
	})
}

// BookEvent handles creating a new booking
func BookCard(c *gin.Context) {
	var input models.Booking
	// 1. Parse the input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a unique user ID using UUID
	booking_id, err := uuid.New()
	if err != nil {
		fmt.Println("Error generating UUID:", err)
		return
	}

	// 3. Create a new Booking
	booking := models.Booking{
		ID:       booking_id.String(),
		EventID:  input.EventID,
		UserID:   input.UserID, // Card creator's user ID
		BookedBy: input.BookedBy,
		Title:    input.Title,
		Status:   input.Status,
	}

	// 4. Save the booking to the database
	if err := config.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	// 5. Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Event booked successfully",
		"booking": booking,
	})
}

// RetrieveMyBookedCards handles fetching all booked cards for a specific user
func RetrieveMyBookedCards(c *gin.Context) {
	// 1. Extract 'booked_by' from query parameters
	bookedBy := c.Query("booked_by")

	// Validate the input
	if bookedBy == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID (booked_by) is required"})
		return
	}

	// 2. Query the database to fetch bookings for the given 'booked_by' ID
	var bookings []models.Booking
	if err := config.DB.Where("booked_by = ?", bookedBy).Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	// 3. Return success response with the list of bookings
	c.JSON(http.StatusOK, gin.H{
		"message":  "Bookings retrieved successfully",
		"bookings": bookings,
	})
}

// RetrieveMyBookedCards handles fetching all booked cards for a specific user
func RetrieveMyBookedCardsRequestToTalent(c *gin.Context) {
	// 1. Extract 'booked_by' from query parameters
	talentId := c.Query("userId")

	// Validate the input
	if talentId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID (card creator) is required"})
		return
	}

	// 2. Query the database to fetch bookings for the given 'booked_by' ID
	var bookings []models.Booking
	if err := config.DB.Where("user_id = ?", talentId).Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	// 3. Return success response with the list of bookings
	c.JSON(http.StatusOK, gin.H{
		"message":  "Bookings retrieved successfully",
		"bookings": bookings,
	})
}

func HandleUpdateBookingStatus(c *gin.Context) {
	// Extract booking ID and user ID from query parameters
	bookingID := c.Query("bookingId")
	userID := c.Query("user_id")
	println(bookingID)
	println(userID)

	// Validate the booking ID and user ID
	if bookingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking ID is required"})
		return
	}
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Input for status update
	type StatusUpdateInput struct {
		Status string `json:"status" binding:"required"`
	}

	var input StatusUpdateInput

	// Bind JSON payload
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: 'status' is required"})
		return
	}

	// Validate status value
	if input.Status != "Accepted" && input.Status != "Declined" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Allowed values are 'Accepted' or 'Declined'"})
		return
	}

	// Fetch and update the booking status
	var booking models.Booking
	if err := config.DB.First(&booking, "id = ?", bookingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	// Update the status
	booking.Status = input.Status
	if err := config.DB.Save(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking status updated successfully",
		"booking": booking,
	})
}
