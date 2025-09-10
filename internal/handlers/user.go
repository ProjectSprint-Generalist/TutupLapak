package handlers

import (
	"net/http"

	"tutuplapak/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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
			Error:   "Server error",
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
