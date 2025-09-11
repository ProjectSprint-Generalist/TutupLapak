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
	FileID           uint            `json:"fileId" gorm:"not null"`
	FileURI          string          `json:"fileUri" gorm:"type:text"`
	FileThumbnailURI string          `json:"fileThumbnailUri" gorm:"type:text"`
	CreatedAt        time.Time       `json:"createdAt"`
	UpdatedAt        time.Time       `json:"updatedAt"`
}
