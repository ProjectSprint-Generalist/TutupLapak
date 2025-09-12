package models

import "time"

type Purchase struct {
	ID                  string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	SenderName          string         `json:"senderName" gorm:"not null"`
	SenderContactType   ContactType    `json:"senderContactType" gorm:"not null"`
	SenderContactDetail string         `json:"senderContactDetail" gorm:"not null"`
	TotalPrice          uint           `json:"totalPrice" gorm:"not null"`
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
	PurchaseItems       []PurchaseItem `json:"purchaseItems" gorm:"foreignKey:PurchaseID;references:ID"`
}

type PurchaseItem struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	PurchaseID string    `json:"purchaseId" gorm:"not null;type:uuid"`
	ProductID  uint      `json:"productId" gorm:"not null"`
	Quantity   uint      `json:"quantity" gorm:"not null"`
	Price      uint      `json:"price" gorm:"not null"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	Product    Product   `json:"product" gorm:"foreignKey:ProductID"`
}

type PurchasePaymentProof struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	PurchaseID   string    `json:"purchaseId" gorm:"not null;type:uuid"`
	FileUploadID uint      `json:"fileUploadId" gorm:"not null"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type PurchasedItems struct {
	ProductID string `json:"productId" binding:"required"`
	Quantity  uint   `json:"qty" binding:"required,min=2"`
}

type PurchaseRequest struct {
	PurchasedItems      []PurchasedItems `json:"purchasedItems" binding:"required,min=1,dive,required"`
	SenderName          string           `json:"senderName" binding:"required,min=4,max=55"`
	SenderContactType   ContactType      `json:"senderContactType" binding:"required,oneof=phone email"`
	SenderContactDetail string           `json:"senderContactDetail" binding:"required"`
}

type ProcessPurchaseRequest struct {
	FileIDs []string `json:"fileIds" binding:"required,min=1,dive,required"`
}

type SellerPaymentInfo struct {
	BankAccountName   string `json:"bankAccountName"`
	BankAccountHolder string `json:"bankAccountHolder"`
	BankAccountNumber string `json:"bankAccountNumber"`
	TotalPrice        uint   `json:"totalPrice"`
}

type PurchaseResponse struct {
	PurchaseID     string                  `json:"purchaseId"`
	PurchasedItems []PurchasedItemResponse `json:"purchasedItems"`
	TotalPrice     uint                    `json:"totalPrice"`
	PaymentDetails []SellerPaymentInfo     `json:"paymentDetails"`
}

type PurchasedItemResponse struct {
	ProductID        string `json:"productId"`
	Name             string `json:"name"`
	Category         string `json:"category"`
	Qty              uint   `json:"qty"`
	Price            uint   `json:"price"`
	SKU              string `json:"sku"`
	FileID           string `json:"fileId"`
	FileURI          string `json:"fileUri"`
	FileThumbnailURI string `json:"fileThumbnailUri"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}
