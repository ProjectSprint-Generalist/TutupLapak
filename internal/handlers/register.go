package handlers

import (
	"errors"
	"net/http"
	"strings"
	"tutuplapak/internal/middleware"
	"tutuplapak/internal/models"
	"tutuplapak/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type RegisterHandler struct {
	db *gorm.DB
}

// NewRegisterHandler initializes RegisterHandler with the given DB
func NewRegisterHandler(db *gorm.DB) *RegisterHandler {
	return &RegisterHandler{
		db: db,
	}
}

func (h *RegisterHandler) RegisterEmail(context *gin.Context) {
	var inputUser models.LoginEmailInput
	if err := context.ShouldBindJSON(&inputUser); err != nil {
		response := models.ErrorResponse{
			Success: false,
			Error:   "Validation error",
			Code:    http.StatusBadRequest,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}

	// Validate email and password input
	if err := utils.Validate(&inputUser); err != nil {
		response := models.ErrorResponse{
			Success: false,
			Error:   "Validation error",
			Code:    http.StatusBadRequest,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(inputUser.Password)
	if err != nil {
		response := models.ErrorResponse{
			Success: false,
			Error:   "Internal server error",
			Code:    http.StatusInternalServerError,
		}
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	// Create a new user
	user := &models.User{
		Email:    inputUser.Email,
		Password: hashedPassword,
	}

	// Check duplicate email
	var existing models.User
	if err := h.db.Where("LOWER(email) = ?", strings.ToLower(inputUser.Email)).First(&existing).Error; err == nil {
		response := models.ErrorResponse{
			Success: false,
			Error:   "Email already registered",
			Code:    http.StatusConflict,
		}
		context.JSON(http.StatusConflict, response)
		return
	}

	if err := h.db.Create(&user).Error; err != nil {
		response := models.ErrorResponse{
			Success: false,
			Error:   "Failed to create user",
			Code:    http.StatusInternalServerError,
		}
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	// Generate JWT Token
	token, err := middleware.GenerateToken(user)
	if err != nil {
		response := models.ErrorResponse{
			Success: false,
			Error:   "Internal server error",
			Code:    http.StatusInternalServerError,
		}
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"email": user.Email,
		"phone": "", // empty string if first registering
		"token": token,
	})
}

func (h *RegisterHandler) RegisterPhone(context *gin.Context) {
	var inputUser models.PhoneUser

	// Bind JSON to Register
	if err := context.ShouldBindJSON(&inputUser); err != nil {
		response := models.ErrorResponse{
			Success: false,
			Error:   "Validation Error",
			Code:    http.StatusBadRequest,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}

	// Validate Phone Number
	if err := utils.PhoneValidation(inputUser.Phone); err != nil {
		response := models.ErrorResponse{
			Success: false,
			Error:   "Validation Error",
			Code:    http.StatusBadRequest,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}

	// Validate Password
	if err := utils.PasswordValidation(inputUser.Password); err != nil {
		response := models.ErrorResponse{
			Success: false,
			Error:   "Validation Error",
			Code:    http.StatusBadRequest,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(inputUser.Password)
	if err != nil {
		response := models.ErrorResponse{
			Success: false,
			Error:   "Internal server error",
			Code:    http.StatusInternalServerError,
		}
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	// Create a new user
	user := &models.User{
		Phone:    inputUser.Phone,
		Password: hashedPassword,
	}

	if err := h.db.Create(&user).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			context.JSON(http.StatusConflict, models.ErrorResponse{
				Success: false,
				Error:   "Phone number already registered",
				Code:    http.StatusConflict,
			})
			return
		}

		context.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to create user",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Generate JWT Token
	token, err := middleware.GenerateToken(user)
	if err != nil {
		response := models.ErrorResponse{
			Success: false,
			Error:   "Internal server error",
			Code:    http.StatusInternalServerError,
		}
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"phone": user.Phone,
		"email": "",
		"token": token,
	})
}
