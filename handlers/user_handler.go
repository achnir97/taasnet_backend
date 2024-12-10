package handlers

import (
	"net/http"

	"taas-api/config"
	"taas-api/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Signup Handler
func Signup(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	user := models.User{
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Login Handler
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
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
		UserID:         input.UserID,
		Title:          input.Title,
		Description:    input.Description,
		Category:       input.Category,
		EventType:      input.EventType,
		Price:          input.Price,
		EventDate:      input.EventDate, // Already a time.Time value
		EventTime:      input.EventTime, // Already a time.Time value
		Participants:   input.Participants,
		VideoURL:       input.VideoURL,
		AvailableTimes: input.AvailableTimes,
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

// BookEvent handles creating a new booking
func BookEvent(c *gin.Context) {
	var input models.Booking
	// 1. Parse the input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Find the Card (Event) by EventID
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

// SaveVideoControl handles saving video control states
func SaveVideoControl(c *gin.Context) {
	var input models.VideoControl

	// Parse the input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create or update video control data
	videoControl := models.VideoControl{
		VideoURL:   input.VideoURL,
		Action:     input.Action,
		StartTime:  input.StartTime,
		EndTime:    input.EndTime,
		PausedTime: input.PausedTime,
	}

	// Save to database
	if err := config.DB.Save(&videoControl).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save video control"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message":       "Video control updated successfully",
		"video_control": videoControl,
	})
}
