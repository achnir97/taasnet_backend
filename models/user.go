package models

type User struct {
	ID       string `json:"id" gorm:"primaryKey"` // UUID as primary key
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
}
