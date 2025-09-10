package handlers

import (
	"net/http"

	"tutuplapak/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

// NewUserHandler creates a new user handler with dependency injection
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// GetUser returns current user profile (GET /v1/user)
func (h *UserHandler) GetUser(c *gin.Context) {
	// Extract user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	// Fetch user from database
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Success: false,
				Error:   "User not found",
				Code:    http.StatusNotFound,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Server Error",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Map to response model with only needed fields
	userResponse := models.UserResponse{
		Email:             user.Email,
		Phone:             user.Phone,
		FileID:            user.FileID,
		FileURI:           user.FileURI,
		FileThumbnailURI:  user.FileThumbnailURI,
		BankAccountName:   user.BankAccountName,
		BankAccountHolder: user.BankAccountHolder,
		BankAccountNumber: user.BankAccountNumber,
	}

	c.JSON(http.StatusOK, userResponse)
}

// UpdateUser updates user profile (PUT /v1/user)
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Expired / invalid / missing request token",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Validation error",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Success: false,
				Error:   "User not found",
				Code:    http.StatusNotFound,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Server Error",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	if req.FileID != "" {
		var fileUpload models.FileUpload
		if err := h.db.Where("id = ? AND user_id = ?", req.FileID, userID).First(&fileUpload).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusBadRequest, models.ErrorResponse{
					Success: false,
					Error:   "fileId is not valid / exists",
					Code:    http.StatusBadRequest,
				})
				return
			}

			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Error:   "Server Error",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		user.FileID = req.FileID
		user.FileURI = fileUpload.FileURI
		user.FileThumbnailURI = fileUpload.FileURI
	}

	user.BankAccountName = req.BankAccountName
	user.BankAccountHolder = req.BankAccountHolder
	user.BankAccountNumber = req.BankAccountNumber

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Server Error",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	userResponse := models.UserResponse{
		Email:             user.Email,
		Phone:             user.Phone,
		FileID:            user.FileID,
		FileURI:           user.FileURI,
		FileThumbnailURI:  user.FileThumbnailURI,
		BankAccountName:   user.BankAccountName,
		BankAccountHolder: user.BankAccountHolder,
		BankAccountNumber: user.BankAccountNumber,
	}

	c.JSON(http.StatusOK, userResponse)
}
