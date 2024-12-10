package handlers

import (
	"net/http"

	"taas-api/config"
	"taas-api/models"

	"github.com/gin-gonic/gin"
)

func CreateCard(c *gin.Context) {
	var input models.Card

	// Parse the request body
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the card
	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create card"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card created successfully", "card": input})
}
