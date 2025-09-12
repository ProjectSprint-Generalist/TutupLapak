package models

import "time"

type FileUpload struct {
	FileID           string    `json:"file_id" gorm:"primaryKey;not null;uniqueIndex"`
	FileName         string    `json:"file_name" gorm:"not null"`
	FileSize         int64     `json:"file_size" gorm:"not null"`
	FileType         string    `json:"file_type" gorm:"not null"`
	FileURI          string    `json:"file_uri" gorm:"not null;unique"`
	FileThumbnailURI string    `json:"file_thumbnail_uri" gorm:"not null"`
	UserID           *uint     `json:"user_id" gorm:"default:null"`
	User             User      `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type FileUploadResponse struct {
	FileID           string `json:"fileId"`
	FileURI          string `json:"fileUri"`
	FileThumbnailURI string `json:"fileThumbnailUri"`
}
