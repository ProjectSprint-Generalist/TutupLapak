package handlers

import (
	"net/http"

	"tutuplapak/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (h *UserHandler) LinkEmail(c *gin.Context) {
	var payload models.LinkEmailRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		response := models.ErrorResponse{
			Success: false,
			Error:   "Invalid input: please provide a valid email request",
			Code:    http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := utils.EmailValidation(payload.Email); err != nil {
		response := models.ErrorResponse{
			Success: false,
			Error:   "Invalid input: invalid email format",
			Code:    http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		fmt.Println("user_id not found in context")
		return
	}

	var user models.User

	if err := h.db.Model(&user).
		Where("id = ?", userID).
		Update("email", payload.Email).Error; err != nil {

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Success: false,
				Error:   "Email already registered",
				Code:    http.StatusConflict,
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

	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to fetch updated user",
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

// LinkPhone (POST /v1/user/link/phone)
func (h *UserHandler) LinkPhone(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Expired / invalid / missing request token",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	var req models.LinkPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Validation error",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := utils.PhoneValidation(req.Phone); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid phone number format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var user models.User
	if err := h.db.Model(&models.User{}).
		Clauses(clause.Returning{}).
		Where("id = ?", userID).
		Update("phone", req.Phone).
		Scan(&user).Error; err != nil {

		// Handle nomor sudah dipakai user lain
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Success: false,
				Error:   "Phone number already linked to another account",
				Code:    http.StatusConflict,
			})
			return
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
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
