package models

import (
	"time"
)

type ProductCategory string

const (
	Food      ProductCategory = "Food"
	Beverage  ProductCategory = "Beverage"
	Clothes   ProductCategory = "Clothes"
	Furniture ProductCategory = "Furniture"
	Tools     ProductCategory = "Tools"
)

// Data Transfer Object Login Input
type LoginEmailInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

type PhoneUser struct {
	Phone    string `json:"phone" binding:"required"`
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

type ProductInput struct {
	Name     string          `json:"name" binding:"required,min=4,max=32"`
	Category ProductCategory `json:"category" binding:"required,oneof=Food Beverage Clothes Furniture Tools"` // Validate both at compile time and runtime
	Qty      uint            `json:"qty" binding:"required,min=1"`
	Price    uint            `json:"price" binding:"required,min=100"`
	SKU      string          `json:"sku" binding:"required,max=32"`
	FileID   uint          `json:"fileId" binding:"required,min=1"` // Should be a valid fileId (received from file upload endpoint), check at runtime
}

type ProductOutput struct {
	ProductID        string      `json:"productId"`
	Name             string    `json:"name"`
	Category         string    `json:"category"`
	Quantity         uint      `json:"quantity"`
	Price            uint      `json:"price"`
	SKU              string    `json:"sku"`
	FileID           uint      `json:"fileId"`
	FileURI          string    `json:"fileUri"`
	FileThumbnailURI string    `json:"fileThumbnailUri"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}
