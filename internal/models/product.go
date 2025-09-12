package models

import "time"

type Product struct {
	ID               uint            `json:"productId" gorm:"primaryKey"`
	UserID           uint            `json:"-" gorm:"index;not null"`
	Name             string          `json:"name" gorm:"type:varchar(32);not null"`
	Category         ProductCategory `json:"category" gorm:"type:varchar(16);not null"`
	Qty              uint            `json:"qty" gorm:"not null;check:qty > 0"`
	Price            uint            `json:"price" gorm:"not null;check:price >= 100"`
	SKU              string          `json:"sku" gorm:"type:varchar(32);not null"`
	FileID           string          `json:"fileId" gorm:"not null"`
	FileURI          string          `json:"fileUri" gorm:"type:text"`
	FileThumbnailURI string          `json:"fileThumbnailUri" gorm:"type:text"`
	CreatedAt        time.Time       `json:"createdAt"`
	UpdatedAt        time.Time       `json:"updatedAt"`
}

// Request payload for update
type UpdateProductRequest struct {
	Name     string `json:"name" binding:"required,min=4,max=32"`
	Category string `json:"category" binding:"required,oneof=Food Beverage Clothes Furniture Tools"`
	Qty      uint   `json:"qty" binding:"required,min=1"`
	Price    uint   `json:"price" binding:"required,min=100"`
	SKU      string `json:"sku" binding:"required,min=1,max=32"`
	FileID   string `json:"fileId" binding:"required"`
}

// Response payload
type ProductResponse struct {
	ProductID        string    `json:"productId"`
	Name             string    `json:"name"`
	Category         string    `json:"category"`
	Qty              uint      `json:"qty"`
	Price            uint      `json:"price"`
	SKU              string    `json:"sku"`
	FileID           string    `json:"fileId"`
	FileURI          string    `json:"fileUri"`
	FileThumbnailURI string    `json:"fileThumbnailUri"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}
