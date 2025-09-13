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

type ContactType string

const (
	ContactTypePhone ContactType = "phone"
	ContactTypeEmail ContactType = "email"
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

type LinkPhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
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
	FileID   string            `json:"fileId" binding:"required,min=1"`
}

type ProductOutput struct {
	ProductID        string    `json:"productId"`
	Name             string    `json:"name"`
	Category         string    `json:"category"`
	Quantity         uint      `json:"quantity"`
	Price            uint      `json:"price"`
	SKU              string    `json:"sku"`
	FileID           string      `json:"fileId"`
	FileURI          string    `json:"fileUri"`
	FileThumbnailURI string    `json:"fileThumbnailUri"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type ProductQueryParams struct {
	Limit     int             `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset    int             `form:"offset" binding:"omitempty,min=0"`
	ProductID string          `form:"productId" binding:"omitempty"`
	SKU       string          `form:"sku" binding:"omitempty"`
	Category  ProductCategory `form:"category" binding:"omitempty,oneof=Food Beverage Clothes Furniture Tools"`
	SortBy    string          `form:"sortBy" binding:"omitempty,oneof=newest oldest cheapest expensive"`
}

type ProductListResponse struct {
	Success bool            `json:"success"`
	Data    []ProductOutput `json:"data"`
	Total   int64           `json:"total"`
	Limit   int             `json:"limit"`
	Offset  int             `json:"offset"`
}
