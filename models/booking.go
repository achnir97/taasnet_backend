package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type AvailableTimeSlots struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	TalentID       string         `json:"talent_id" binding:"required"`
	AvailableDate  time.Time      `json:"available_date" binding:"required"`
	AvailableSlots TimeRanges     `gorm:"type:jsonb" json:"available_slots" binding:"required"` // JSON format
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
type TimeRange struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type TimeRanges []TimeRange
type SessionType string

const (
	CoffeeCall SessionType = "CoffeeCall"
	Regular    SessionType = "Regular"
)

type BookingStatus string

const (
	Scheduled BookingStatus = "Scheduled"
	Completed BookingStatus = "Completed"
	Cancelled BookingStatus = "Cancelled"
)

type PaymentStatus string

const (
	Paid    PaymentStatus = "Paid"
	Pending PaymentStatus = "Pending"
)

type BookingRequests struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	BookingID       string         `json:"booking_id" binding:"required"`
	CardID          string         `json:"card_id" binding:"required"`
	CardTitle       string         `json:"card_title" binding:"required"`
	UserID          string         `json:"user_id" binding:"required"`
	TalentID        string         `json:"talent_id" binding:"required"`
	SessionType     SessionType    `gorm:"type:text" json:"session_type" binding:"required"`
	BookedTime      TimeRanges     `gorm:"type:jsonb" json:"booked_time" binding:"required"` // combined start and end time
	BookingDate     time.Time      `json:"booking_date" binding:"required"`
	Status          BookingStatus  `gorm:"type:text" json:"status" binding:"required"`
	PaymentStatus   PaymentStatus  `gorm:"type:text" json:"payment_status" binding:"required"`
	SpecialRequests string         `json:"special_requests"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// Value implements the driver.Valuer interface.
func (t TimeRanges) Value() (driver.Value, error) {
	if len(t) == 0 {
		return nil, nil
	}
	j, err := json.Marshal(t)
	return j, err
}

// Scan implements the sql.Scanner interface.
func (t *TimeRanges) Scan(value interface{}) error {
	if value == nil {
		*t = []TimeRange{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	var decoded []TimeRange
	err := json.Unmarshal(bytes, &decoded)
	if err != nil {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", string(bytes), "Error: ", err))
	}
	*t = decoded
	return nil
}
