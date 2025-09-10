package handlers

import "gorm.io/gorm"

// LoginHandler manage user login
type LoginHandler struct {
	db *gorm.DB
}

// RegisterHandler manage user registration
type RegisterHandler struct {
	db *gorm.DB
}

type UserHandler struct {
	db *gorm.DB
}

// NewUserHandler creates a new user handler with dependency injection
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}