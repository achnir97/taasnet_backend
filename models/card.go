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

type Booking struct {
	ID        uint      `gorm:"primaryKey" json:"id"`      // Use string for UUID
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
