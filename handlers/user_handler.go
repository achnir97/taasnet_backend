package handlers

import (
	"fmt"
	"net/http"
	"taas-api/config"
	"taas-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"time"

	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"storj.io/common/uuid"
)

// Initialize the validator globally
var validate *validator.Validate = validator.New()

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

// RegisterUser handles user registration requests.
func RegisterUser(c *gin.Context) {

	var input struct { // Primary Key
		FirstName   string         `gorm:"size:100;not null" json:"first_name"`                 // User's first name
		LastName    string         `gorm:"size:100;not null" json:"last_name"`                  // User's last name
		Phone       string         `gorm:"size:20;unique;not null" json:"phone"`                // User's phone number
		Email       string         `gorm:"size:100;unique;not null" json:"email"`               // User's email address
		Password    string         `gorm:"size:255;not null" json:"password"`                   // Hashed password
		AccountType string         `gorm:"size:50;not null;default:'User'" json:"account_type"` // Enum: "User" or "Talent"
		CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`                    // User creation timestamp
		UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`                    // User profile last updated timestamp
		DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`                   // Soft delete
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Check if email already exists
	var existingUser models.Users_ref
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}
	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	input.Password = string(hashedPassword)

	//Generate a unique user ID using UUID
	uniqueUserID, err := uuid.New()
	if err != nil {
		fmt.Println("Error generating UUID:", err)
		return
	}
	// Create the user
	user := models.Users_ref{
		UserID:      uniqueUserID.String(), // Convert to UUI
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Phone:       input.Phone,
		AccountType: input.AccountType,
		Email:       input.Email,
		Password:    string(hashedPassword),
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
func SignIn(c *gin.Context) {
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
	var user models.Users_ref
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
		"userId":  user.UserID, // Include user ID in the response
	})
}

// RegisterTalent handles the talent registration
func RegisterTalent(c *gin.Context) {
	// Parse incoming JSON data into a struct
	var talent struct {
		UserID          string   `json:"user_id" binding:"required"`
		TalentName      string   `json:"talent_name" binding:"required"`
		Category        string   `json:"category" binding:"required"`
		Bio             string   `json:"bio" binding:"required"`
		PortfolioURL    string   `json:"portfolio_url"`
		Skills          []string `json:"skills" binding:"required"` // Array of skills from the request
		ProfileImage    string   `json:"profile_image"`             // Base64 encoded or file path
		ExperienceLevel string   `json:"experience_level" binding:"required"`
	}

	// Bind JSON data and validate
	if err := c.ShouldBindJSON(&talent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data: " + err.Error()})
		return
	}

	// Join skills into a single comma-separated string
	skillsString := strings.Join(talent.Skills, ",")

	// Handle profile image upload
	imagePath, err := handleImageUpload(talent.ProfileImage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save profile image"})
		return
	}

	//Generate a unique user ID using UUID
	uniquetalentid, err := uuid.New()
	if err != nil {
		fmt.Println("Error generating UUID:", err)
		return
	}
	// Create TalentRegistration object for saving to DB
	talentRecord := models.TalentRegistration{
		UserID:          talent.UserID,
		TalentID:        uniquetalentid.String(),
		TalentName:      talent.TalentName,
		Category:        talent.Category,
		Bio:             talent.Bio,
		ExperienceLevel: talent.ExperienceLevel,
		PortfolioLink:   talent.PortfolioURL,
		Skills:          skillsString,
		ProfileImageURL: imagePath,
	}

	// Save to database
	if err := config.DB.Create(&talentRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save talent data"})
		return
	}

	// Success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Talent registered successfully!",
		"skills":  skillsString,
	})
}

func GetTalentAccounts(c *gin.Context) {
	userID := c.Query("user_id") // Get the user ID from the request parameters

	var talents []models.TalentRegistration // Slice to store the talents associated with the user

	// Query the database for all talent accounts associated with the user ID
	if err := config.DB.Where("user_id = ?", userID).Find(&talents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch talents"})
		return
	}

	// Check if no talents exist for the user
	if len(talents) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No talents created yet"})
		return
	}

	// Prepare the response data
	talentData := make([]map[string]string, len(talents))
	for i, talent := range talents {
		talentData[i] = map[string]string{
			"talent_id":   talent.TalentID,
			"talent_name": talent.TalentName,
		}
	}

	// Success response with the list of TalentIDs and TalentNames
	c.JSON(http.StatusOK, gin.H{
		"message":     "Talent accounts fetched successfully",
		"user_id":     userID,
		"talent_data": talentData,
	})
}

// handleImageUpload saves the profile image to a folder and returns the file path
func handleImageUpload(base64Image string) (string, error) {
	// For now, let's save the image as a static file (you can implement base64 decoding if needed)
	if base64Image == "" {
		return "", nil
	}

	// Static folder for uploaded images
	uploadDir := "uploads/"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	// Generate a unique filename
	imageName := fmt.Sprintf("profile_%d.png", os.Getpid()) // Unique filename
	imagePath := filepath.Join(uploadDir, imageName)

	// Mock saving the image (for real base64 decoding, use a library)
	file, err := os.Create(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Write image data (placeholder, real implementation may differ)
	_, err = file.WriteString("mock_image_data")
	if err != nil {
		return "", err
	}

	return imagePath, nil
}
