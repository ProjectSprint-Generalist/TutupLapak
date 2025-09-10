package models

import "time"

// Full Product model (DB entity)
type Product struct {
	BaseEntity
	Name             string    `gorm:"column:name;type:varchar(32);not null" json:"name"`
	Category         string    `gorm:"column:category;type:varchar(32);not null" json:"category"`
	Qty              int       `gorm:"column:qty;not null" json:"qty"`
	Price            float64   `gorm:"column:price;type:decimal(10,2);not null" json:"price"`
	Sku              string    `gorm:"column:sku;type:varchar(32);not null" json:"sku"`
	FileID           string    `gorm:"column:file_id;type:varchar(100);not null" json:"fileId"`
	FileURI          string    `gorm:"column:file_uri;type:text;not null" json:"fileUri"`
	FileThumbnailURI string    `gorm:"column:file_thumbnail_uri;type:text;not null" json:"fileThumbnailUri"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// Request payload for update
type UpdateProductRequest struct {
	Name     string  `json:"name" binding:"required,min=4,max=32"`
	Category string  `json:"category" binding:"required,oneof=Food Beverage Clothes Furniture Tools"`
	Qty      int     `json:"qty" binding:"required,min=1"`
	Price    float64 `json:"price" binding:"required,min=100"`
	Sku      string  `json:"sku" binding:"required,min=0,max=32"`
	FileId   string  `json:"fileId" binding:"required"`
}

// Response payload
type ProductResponse struct {
	ProductID        string    `json:"productId"`
	Name             string    `json:"name"`
	Category         string    `json:"category"`
	Qty              int       `json:"qty"`
	Price            float64   `json:"price"`
	Sku              string    `json:"sku"`
	FileID           string    `json:"fileId"`
	FileURI          string    `json:"fileUri"`
	FileThumbnailURI string    `json:"fileThumbnailUri"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}