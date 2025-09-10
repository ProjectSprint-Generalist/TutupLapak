package models

// Data Transfer Object Login Input
type LoginEmailInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

type LoginPhoneInput struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

type LoginPhoneOutput struct {
	Phone string `json:"phone"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type LinkEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}
