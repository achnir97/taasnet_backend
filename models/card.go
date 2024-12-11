package models

import (
	"time"
)

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

// Booking model for event bookings
type Booking struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	EventID   string    `gorm:"not null" json:"event_id"`         // Refers to Card ID
	UserID    string    `gorm:"not null" json:"user_id"`          // ID of the Card creator
	BookedBy  uint      `gorm:"not null" json:"booked_by"`        // ID of the user making the booking
	Title     string    `gorm:"not null" json:"title"`            // Title of the event
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // Booking creation time
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
