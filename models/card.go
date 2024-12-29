package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type StringSlice []string

// Implement Scan for reading from the database
func (s *StringSlice) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), s)
}

// Implement Value for writing to the database
func (s StringSlice) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// TalentRegistration represents the structure for talent registration data
type TalentRegistration struct {
	UserID          string         `gorm:"not nul" json:"user_id"`                   // Primary Keyst
	TalentID        string         `gorm:"not null" json:"talent_id"`                // Auto-increment Primary Key
	TalentName      string         `gorm:"size:255;not null" json:"talent_name"`     // Name of the talent
	Category        string         `gorm:"size:100;not null" json:"category"`        // Talent category
	Bio             string         `gorm:"type:text;not null" json:"bio"`            // Short bio
	Skills          string         `gorm:"type:text" json:"skills"`                  // Array of skills
	PortfolioLink   string         `gorm:"size:255" json:"portfolio_link"`           // Portfolio URL
	ProfileImageURL string         `gorm:"size:255" json:"profile_image_url"`        // Uploaded image URL
	ExperienceLevel string         `gorm:"size:50;not null" json:"experience_level"` // Experience Level
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`         // Timestamp when created
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`         // Timestamp when updated
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`        // Soft delete
}

// ServiceCard represents the individual service offerings by talents
type ServiceCard struct {
	CardID          string         `gorm:"not null;" json:"card_id"`                   // Primary Key
	TalentID        string         `gorm:"not null;index" json:"talent_id"`            // Foreign Key: Links to Talent table
	CardTitle       string         `gorm:"size:255;not null" json:"card_title"`        // Title of the service
	CardDescription string         `gorm:"type:text;not null" json:"card_description"` // Detailed description
	Suit            string         `gorm:"size:50;not null" json:"suit"`               // Enum: "Heart", "Spade", "Diamond", "Clover"
	Price           int            `gorm:"not null" json:"price"`                      // Price for the service
	Duration        int            `gorm:"not null" json:"duration"`                   // Duration in minutes
	Tags            string         `gorm:"size:255" json:"tags"`                       // Comma-separated tags
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`           // Timestamp for creation
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`           // Timestamp for update
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`          // Soft delete field
}

// Booking represents session bookings for CoffeeCalls and other services
type Bookings struct {
	BookingID       string         `gorm:"primaryKey;autoIncrement" json:"booking_id"` // Primary Key
	CardID          string         `gorm:"not null;index" json:"card_id"`              // Foreign Key: Links to TaaS card
	UserID          string         `gorm:"not null;index" json:"user_id"`              // Foreign Key: Links to User table
	TalentID        string         `gorm:"not null;index" json:"talent_id"`            // Foreign Key: Links to Talent table
	SessionType     string         `gorm:"size:50;not null" json:"session_type"`       // Enum: "CoffeeCall", "Regular"
	StartTime       time.Time      `gorm:"not null" json:"start_time"`                 // Start time of the session
	EndTime         time.Time      `gorm:"not null" json:"end_time"`                   // End time of the session
	Status          string         `gorm:"size:50;not null" json:"status"`             // Enum: "Scheduled", "Completed", "Cancelled"
	PaymentStatus   string         `gorm:"size:50;not null" json:"payment_status"`     // Enum: "Paid", "Pending"
	SpecialRequests string         `gorm:"type:text" json:"special_requests"`          // Special requests made during booking
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`           // Timestamp: Booking created
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`           // Timestamp: Last update
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`          // Soft delete
}

// Session represents the session table in the database
type Session struct {
	SessionID     uint           `gorm:"primaryKey;autoIncrement" json:"session_id"`                                       // Primary Key
	CardID        uint           `gorm:"not null" json:"card_id"`                                                          // Foreign Key: Links to the TaaS card
	UserID        uint           `gorm:"not null" json:"user_id"`                                                          // Foreign Key: Links to the user who booked
	TalentID      uint           `gorm:"not null" json:"talent_id"`                                                        // Foreign Key: Links to the talent providing the session
	SessionType   string         `gorm:"type:enum('CoffeeCall','Regular');default:'Regular'" json:"session_type"`          // Enum: "CoffeeCall" or "Regular"
	StartTime     time.Time      `gorm:"not null" json:"start_time"`                                                       // Scheduled start time
	EndTime       time.Time      `gorm:"not null" json:"end_time"`                                                         // Scheduled end time
	Status        string         `gorm:"type:enum('Scheduled','Completed','Cancelled');default:'Scheduled'" json:"status"` // Enum for session status
	PaymentStatus string         `gorm:"type:enum('Paid','Pending');default:'Pending'" json:"payment_status"`              // Enum for payment status
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`                                                 // Timestamp for when the session was created
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`                                                 // Timestamp for updates
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`                                                // Soft delete (optional)
}

// Message represents the message table in the database
type Message struct {
	MessageID      uint           `gorm:"primaryKey;autoIncrement" json:"message_id"`                             // Primary Key
	CardID         uint           `gorm:"not null" json:"card_id"`                                                // Foreign Key: Links to TaaS card
	UserID         uint           `gorm:"not null" json:"user_id"`                                                // Foreign Key: User who sent the message
	TalentID       uint           `gorm:"not null" json:"talent_id"`                                              // Foreign Key: Talent receiving the message
	MessageContent string         `gorm:"type:text;not null" json:"message_content"`                              // Actual question or request
	MessageType    string         `gorm:"type:enum('Question','Request');default:'Question'" json:"message_type"` // Enum for message type
	BookingID      *uint          `gorm:"default:null" json:"booking_id,omitempty"`                               // Optional FK: Links to a booking
	Status         string         `gorm:"type:enum('Pending','Responded');default:'Pending'" json:"status"`       // Message status
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`                                       // Timestamp for message creation
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`                                       // Timestamp for updates
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`                                      // Optional: Soft delete
}

// Review represents user reviews for completed services
type Review struct {
	ReviewID      uint           `gorm:"primaryKey;autoIncrement" json:"review_id"` // Primary Key
	CardID        uint           `gorm:"not null;index" json:"card_id"`             // Foreign Key: Links to TaaS card
	UserID        uint           `gorm:"not null;index" json:"user_id"`             // Foreign Key: Links to User table
	TalentID      uint           `gorm:"not null;index" json:"talent_id"`           // Foreign Key: Links to Talent table
	Rating        float64        `gorm:"not null" json:"rating"`                    // Numeric rating (e.g., 1-5)
	ReviewContent string         `gorm:"type:text" json:"review_content"`           // Text content of the review
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`          // Timestamp: Review created
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`         // Soft delete
}

// Card model
type Card struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       string    `gorm:"not null"`
	Title        string    `gorm:"not null"`
	Description  string    `gorm:"type:text"`
	Category     string    `gorm:"not null"`
	EventType    string    `gorm:"not null"`
	Price        float64   `gorm:"type:decimal(10,2);default:0"`
	EventDate    time.Time `gorm:"not null"`
	EventTime    time.Time `gorm:"not null"`
	Participants int       `gorm:"default:0"`
	VideoURL     string    `gorm:"type:text"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

type Booking struct {
	ID        string    `json:"id" gorm:"primaryKey"`      // Use string for UUID
	EventID   string    `gorm:"not null" json:"event_id"`  // Event ID (string)
	UserID    string    `gorm:"not null" json:"user_id"`   // User ID (string)
	BookedBy  string    `gorm:"not null" json:"booked_by"` // BookedBy ID (string)
	Title     string    `gorm:"not null" json:"title"`     // Title
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	Status    string    `gorm:"not null" json:"status"`
}

type VideoControl struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	VideoURL   string    `gorm:"type:text;not null" json:"video_url"`
	Action     string    `gorm:"type:varchar(50);not null" json:"action"` // e.g., play, pause, stop
	StartTime  int       `gorm:"not null" json:"start_time"`              // In seconds
	EndTime    int       `gorm:"not null" json:"end_time"`                // In seconds
	PausedTime int       `json:"paused_time"`                             // Last paused time in seconds
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type AvailableTimeSlots struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	TalentID       string         `json:"talent_id" binding:"required"`
	AvailableDate  time.Time      `json:"available_date" binding:"required"`
	AvailableSlots string         `json:"available_slots" binding:"required"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type BookingRequest struct {
	BookingID       string         `json:"booking_id" binding:"required"`
	CardID          string         `json:"card_id" binding:"required"`
	CardTitle       string         `json:"card_title" binding:"required"`
	UserID          string         `json:"user_id" binding:"required"`
	TalentID        string         `json:"talent_id" binding:"required"`
	SessionType     string         `json:"session_type" binding:"required"`   // Enum: CoffeeCall, Regular
	StartTime       time.Time      `json:"start_time" binding:"required"`     // ISO8601 format
	EndTime         time.Time      `json:"end_time" binding:"required"`       // ISO8601 format
	Status          string         `json:"status" binding:"required"`         // Enum: Scheduled, Completed, Cancelled
	PaymentStatus   string         `json:"payment_status" binding:"required"` // Enum: Paid, Pending
	SpecialRequests string         `json:"special_requests"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
