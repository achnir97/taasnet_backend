package handlers

import (
	"fmt"

	"net/http"

	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// Upload directory
const uploadPath = "./uploads"

// FileUploadHandler handles file uploads
func FileUploadHandler(c *gin.Context) {
	// Limit upload size to 10MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20) // 10MB limit

	// Parse the form input
	file, err := c.FormFile("file") // "file" should match the form field name
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}
	// Generate a unique file name using timestamp
	timestamp := time.Now().Unix()
	fileName := fmt.Sprintf("%d_%s", timestamp, filepath.Base(file.Filename))
	filePath := filepath.Join(uploadPath, fileName)

	// Save the uploaded file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}
	// Return file URL
	fileURL := fmt.Sprintf("http://%s/uploads/%s", c.Request.Host, fileName)
	c.JSON(http.StatusCreated, gin.H{"message": "File uploaded successfully", "file_url": fileURL})
}
