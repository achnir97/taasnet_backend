package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       string `json:"id" gorm:"primaryKey"` // UUID as primary key
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
}

// Users represents the user model for registration and management.
type Users_ref struct {
	UserID      string         `gorm:"primaryKey;not null" json:"user_id"`                  // Primary Key
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
