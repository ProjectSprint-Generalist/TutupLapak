package models

// Data Transfer Object Login Input
type LoginEmailInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}
