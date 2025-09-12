package services

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"tutuplapak/internal/config"
	"tutuplapak/internal/models"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/gorm"
)

type MinIOService struct {
	client     *minio.Client
	bucketName string
	db         *gorm.DB
}

func NewMinIOService(cfg *config.Config, db *gorm.DB) (*MinIOService, error) {
	client, err := minio.New(cfg.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIO.AccessKeyID, cfg.MinIO.SecretAccessKey, ""),
		Secure: cfg.MinIO.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	service := &MinIOService{
		client:     client,
		bucketName: cfg.MinIO.BucketName,
		db:         db,
	}

	if err := service.ensureBucketExists(); err != nil {
		return nil, fmt.Errorf("failed to ensure bucket exists: %w", err)
	}

	return service, nil
}

func (s *MinIOService) ensureBucketExists() error {
	ctx := context.Background()
	exists, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		return err
	}

	if !exists {
		err = s.client.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *MinIOService) UploadFile(file *multipart.FileHeader, userID uint) (*models.FileUploadResponse, error) {
	if err := s.validateFile(file); err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Read the entire file into memory for processing
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Generate a unique fileId
	fileID := uuid.New().String()

	// Reset the file reader for upload
	fileReader := bytes.NewReader(fileBytes)

	fileExt := strings.ToLower(filepath.Ext(file.Filename))
	objectName := fmt.Sprintf("uploads/%s%s", fileID, fileExt)

	ctx := context.Background()

	// Upload the original file
	_, err = s.client.PutObject(ctx, s.bucketName, objectName, fileReader, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	fileURI := fmt.Sprintf("%s/%s/%s", s.client.EndpointURL().String(), s.bucketName, objectName)

	// Create thumbnail
	img, _, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Resize to create thumbnail (100x100 maintaining aspect ratio)
	thumbnail := imaging.Resize(img, 100, 100, imaging.Lanczos)

	var thumbnailBuf bytes.Buffer
	var format imaging.Format

	switch strings.ToLower(filepath.Ext(file.Filename)) {
	case ".jpg", ".jpeg":
		format = imaging.JPEG
	case ".png":
		format = imaging.PNG
	default:
		format = imaging.JPEG // Default to JPEG
	}

	err = imaging.Encode(&thumbnailBuf, thumbnail, format)
	if err != nil {
		return nil, fmt.Errorf("failed to encode thumbnail: %w", err)
	}

	// Upload the thumbnail
	thumbnailObjectName := fmt.Sprintf("thumbnails/%s%s", fileID, fileExt)
	_, err = s.client.PutObject(
		ctx,
		s.bucketName,
		thumbnailObjectName,
		bytes.NewReader(thumbnailBuf.Bytes()),
		int64(thumbnailBuf.Len()),
		minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upload thumbnail to MinIO: %w", err)
	}

	thumbnailURI := fmt.Sprintf("%s/%s/%s", s.client.EndpointURL().String(), s.bucketName, thumbnailObjectName)

	fileUpload := &models.FileUpload{
		FileName:         file.Filename,
		FileSize:         file.Size,
		FileType:         file.Header.Get("Content-Type"),
		FileURI:          fileURI,
		FileThumbnailURI: thumbnailURI,
		FileID:           fileID,
	}

	// Only set UserID if it's not zero (for authenticated uploads)
	if userID > 0 {
		fileUpload.UserID = &userID
	}

	if err := s.db.Create(fileUpload).Error; err != nil {
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}

	return &models.FileUploadResponse{
		FileID:           fileID,
		FileURI:          fileURI,
		FileThumbnailURI: thumbnailURI,
	}, nil
}

func (s *MinIOService) validateFile(file *multipart.FileHeader) error {
	const maxSize = 100 * 1024 // 100 KiB
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
	}

	if file.Size > maxSize {
		return fmt.Errorf("file size exceeds 100 KiB limit")
	}

	contentType := file.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return fmt.Errorf("file type %s not allowed. Only JPEG, JPG, and PNG are supported", contentType)
	}

	fileExt := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := []string{".jpg", ".jpeg", ".png"}
	for _, ext := range allowedExts {
		if fileExt == ext {
			return nil
		}
	}

	return fmt.Errorf("file extension %s not allowed. Only .jpg, .jpeg, and .png are supported", fileExt)
}