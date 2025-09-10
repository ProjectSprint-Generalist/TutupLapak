package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	Name              string    `json:"name" gorm:"type:varchar(255)"` 
	Email             string    `json:"email" gorm:"uniqueIndex;not null"`
	Phone             string    `json:"phone" gorm:"type:varchar(20)"`
	Password          string    `json:"-" gorm:"not null"`
	FileID            string    `json:"fileId"`
	FileURI           string    `json:"fileUri"`
	FileThumbnailURI  string    `json:"fileThumbnailUri"`
	BankAccountName   string    `json:"bankAccountName"`
	BankAccountHolder string    `json:"bankAccountHolder"`
	BankAccountNumber string    `json:"bankAccountNumber"`
	ImageURI          string    `json:"imageUri" gorm:"type:text"` 
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// UserResponse represents the response payload for GET /v1/user
type UserResponse struct {
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	FileID            string `json:"fileId"`
	FileURI           string `json:"fileUri"`
	FileThumbnailURI  string `json:"fileThumbnailUri"`
	BankAccountName   string `json:"bankAccountName"`
	BankAccountHolder string `json:"bankAccountHolder"`
	BankAccountNumber string `json:"bankAccountNumber"`
}
