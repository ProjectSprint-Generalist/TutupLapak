package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	Name              string    `json:"name" gorm:"type:varchar(255)"`
	Email             string    `json:"email" gorm:"uniqueIndex;default:null"`
	Phone             string    `json:"phone" gorm:"uniqueIndex;type:varchar(20);default:null"`
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

// InputUser represents the input for user registration
type InputUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

// UpdateUserRequest represents the request payload for PUT /v1/user
type UpdateUserRequest struct {
	FileID            string `json:"fileId" binding:"omitempty"`
	BankAccountName   string `json:"bankAccountName" binding:"required,min=4,max=32"`
	BankAccountHolder string `json:"bankAccountHolder" binding:"required,min=4,max=32"`
	BankAccountNumber string `json:"bankAccountNumber" binding:"required,min=4,max=32"`
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
