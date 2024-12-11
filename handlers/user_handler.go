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

	// Check if email already exists
	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Generate a unique user ID using UUID
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
func BookEvent(c *gin.Context) {
	var input models.Booking
	// 1. Parse the input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Find the Card (Event) by Event ID
	var card models.Card
	if err := config.DB.First(&card, input.EventID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	// 3. Create a new Booking
	booking := models.Booking{
		EventID:  input.EventID,
		UserID:   card.UserID, // Card creator's user ID
		BookedBy: input.BookedBy,
		Title:    card.Title,
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
