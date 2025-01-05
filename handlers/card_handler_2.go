package handlers

import (
	"fmt"
	"net/http"
	"taas-api/config"
	"taas-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// editCard handles editing an existing card
func EditCard(c *gin.Context) {
	fmt.Println("Starting EditCard handler")
	cardId := c.Param("cardId") // Get the card ID from the URL parameter
	fmt.Println("Card ID from URL:", cardId)
	if cardId == "" {
		fmt.Println("Error: Card ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Card ID is required"})
		return
	}
	var card struct {
		// Primary Key
		TalentID        string  `gorm:"not null" json:"talent_id"`         // Foreign Key: Links to Talent table
		CardTitle       *string `json:"card_title"`                        // Title of the service
		CardDescription *string `gorm:"type:text" json:"card_description"` // Detailed description
		Suit            *string `json:"suit"`                              // Enum: "Heart", "Spade", "Diamond", "Clover"
		Price           *int    `    json:"price"`                         // Price for the service
		Duration        *int    `json:"duration"`                          // Duration in minutes
		Tags            *string `json:"tags"`                              // Comma-separated tags
	}
	fmt.Println("Binding Json")
	// Parse the request body into the input struct
	if err := c.ShouldBindJSON(&card); err != nil {
		fmt.Println("Error Binding json:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Inputs from request body:", card)
	fmt.Println("Retrieving existing card from database")

	//Retrieve the existing card from the database
	var existingCard models.ServiceCard
	if err := config.DB.Where("card_id = ?", cardId).First(&existingCard).Error; err != nil {
		fmt.Println("Error Retrieving card from database:", err.Error())
		if err == gorm.ErrRecordNotFound {
			fmt.Println("Error: Card not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve card"})
		return
	}
	fmt.Println("Existing card retrieved:", existingCard)
	//Check if talent id match with the one who created this service card
	fmt.Println("Checking for Talent authorization")

	if card.TalentID != existingCard.TalentID {
		fmt.Println("Error: Unauthorized to edit this card")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to edit this card"})
		return
	}

	fmt.Println("Updating card fields")
	// Create a map of fields to update
	updates := make(map[string]interface{})
	if card.CardTitle != nil {
		updates["card_title"] = *card.CardTitle
	}
	if card.CardDescription != nil {
		updates["card_description"] = *card.CardDescription
	}
	if card.Suit != nil {
		updates["suit"] = *card.Suit
	}
	if card.Price != nil {
		updates["price"] = *card.Price
	}
	if card.Duration != nil {
		updates["duration"] = *card.Duration
	}
	if card.Tags != nil {
		updates["tags"] = *card.Tags
	}

	fmt.Println("Saving updated card to database:", updates)
	// Save the updated card to the database
	if err := config.DB.Model(&models.ServiceCard{}).Where("card_id = ?", cardId).Updates(updates).Error; err != nil {
		fmt.Println("Error saving updated card to database:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update card"})
		return
	}

	var updatedCard models.ServiceCard
	if err := config.DB.Where("card_id = ?", cardId).First(&updatedCard).Error; err != nil {
		fmt.Println("Error retrieving card after update:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve card"})
		return
	}

	fmt.Println("Card updated successfully")

	// Respond with success
	c.JSON(http.StatusOK, gin.H{
		"message": "Card updated successfully",
		"card":    updatedCard, // You can send back the updated data
	})
	fmt.Println("EditCard handler finished")
}

// / deleteCard handles deleting an existing card
func DeleteCard(c *gin.Context) {
	fmt.Println("Starting DeleteCard handler")

	cardId := c.Param("cardId") // Get the card ID from the URL parameter
	fmt.Println("Card ID from URL:", cardId)
	if cardId == "" {
		fmt.Println("Error: Card ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Card ID is required"})
		return
	}
	fmt.Println("Parsing JSON")
	var card struct {
		TalentID string `gorm:"not null" json:"talent_id"`
	}

	// Parse the request body into the input struct
	if err := c.ShouldBindJSON(&card); err != nil {
		fmt.Println("Error binding JSON:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Input from the request:", card)
	fmt.Println("Retrieving card from the database")

	//Retrieve the existing card from the database
	var existingCard models.ServiceCard
	if err := config.DB.Where("card_id = ?", cardId).First(&existingCard).Error; err != nil {
		fmt.Println("Error retrieving card from the database:", err.Error())
		if err == gorm.ErrRecordNotFound {
			fmt.Println("Error: Card not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
			return
		}
		fmt.Println("Error: Failed to retrieve card from database")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve card"})
		return
	}
	fmt.Println("Existing card from database:", existingCard)
	//Check if talent id match with the one who created this service card
	fmt.Println("Checking Talent authorization")
	if card.TalentID != existingCard.TalentID {
		fmt.Println("Error: Unauthorized to delete this card")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to delete this card"})
		return
	}

	fmt.Println("Deleting card from database")

	// Delete the card from the database
	if err := config.DB.Where("card_id = ?", cardId).Delete(&models.ServiceCard{}).Error; err != nil {
		fmt.Println("Error deleting the card from database:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete card"})
		return
	}

	fmt.Println("Card deleted successfully")
	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "Card deleted successfully"})
	fmt.Println("DeleteCard handler finished")
}
