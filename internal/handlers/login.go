package handlers

import (
	"net/http"
	"tutuplapak/internal/middleware"
	"tutuplapak/internal/models"
	"tutuplapak/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewLoginHandler initializes LoginHandler with the given DB
func NewLoginHandler(db *gorm.DB) *LoginHandler {
	return &LoginHandler{
		db: db,
	}
}

func emailValidation(userInput string, userPassword string) *models.ErrorResponse {
	if err := utils.VerifyPassword(userInput, userPassword); err != nil {
		errResponse := &models.ErrorResponse{
			Success: false,
			Error:   "Invalid password",
			Code:    http.StatusUnauthorized,
		}

		return errResponse
	}

	return nil
}

func tokenGeneration(user *models.User) (string, *models.ErrorResponse) {
	token, err := middleware.GenerateToken(user)
	if err != nil {
		response := &models.ErrorResponse{
			Success: false,
			Error:   "Failed to generate token",
			Code:    http.StatusInternalServerError,
		}
		return "", response
	}

	return token, nil

}

// Login handle user login requests
func (h *LoginHandler) LoginEmail(context *gin.Context) {

	var inputUser models.LoginEmailInput

	// Bind JSON
	if err := context.ShouldBindJSON(&inputUser); err != nil {
		response := models.ErrorResponse{
			Success: false,
			Error:   "Invalid input: please provide a valid email and password",
			Code:    http.StatusBadRequest,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}

	//  Check existing user
	var user models.User
	if err := h.db.Where("email = ?", inputUser.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response := models.ErrorResponse{
				Success: false,
				Error:   "User does not exist",
				Code:    http.StatusNotFound,
			}
			context.JSON(http.StatusNotFound, response)
			return
		}
		response := models.ErrorResponse{
			Success: false,
			Error:   "Server error",
			Code:    http.StatusInternalServerError,
		}
		context.JSON(http.StatusInternalServerError, response)
		return
	}

	err := emailValidation(inputUser.Email, inputUser.Password)
	// Verify Password
	if err != nil {
		context.JSON(http.StatusUnauthorized, err)
		return
	}

	// Generate Token
	token, err := tokenGeneration(&user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, err)
	}

	// Login successfully
	context.JSON(http.StatusOK, gin.H{
		"email": user.Email,
		"phone": "", // can be replace with: user.Phone if later exists
		"token": token,
	})
}
