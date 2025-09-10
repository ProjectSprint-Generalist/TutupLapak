package handlers

import (
	"net/http"
	"tutuplapak/internal/middleware"
	"tutuplapak/internal/models"
	"tutuplapak/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginHandler struct {
	db *gorm.DB
}

// NewLoginHandler initializes LoginHandler with the given DB
func NewLoginHandler(db *gorm.DB) *LoginHandler {
	return &LoginHandler{
		db: db,
	}
}

func passwordValidation(userInput string, userPassword string) *models.ErrorResponse {
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

	if err := utils.EmailValidation(inputUser.Email); err != nil {
		response := models.ErrorResponse{
			Success: false,
			Error:   "Invalid input: invalid email format",
			Code:    http.StatusBadRequest,
		}
		context.JSON(http.StatusBadRequest, response)
		return
	}

	if err := utils.PasswordLengthValidation(inputUser.Password); err != nil {
		response := models.ErrorResponse{
			Success: false,
			Error:   "Invalid input: password length must be 8–32 characters",
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

	// Password Hash Validation
	err := passwordValidation(inputUser.Password, user.Password)
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

func (h *LoginHandler) LoginPhone(ctx *gin.Context) {
	var inputUser models.LoginPhoneInput

	// Check if JSON input is valid
	if err := ctx.ShouldBindJSON(&inputUser); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid input: please provide a valid phone and password",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Check if phone number is valid
	if err := utils.PhoneValidation(inputUser.Phone); err != nil {
		ctx.JSON(err.Code, err)
		return
	}

	// Check if password is valid
	if err := utils.PasswordValidation(inputUser.Password); err != nil {
		ctx.JSON(err.Code, err)
		return
	}

	// Check if user exists
	var user models.User
	if err := h.db.Where("phone = ?", inputUser.Phone).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, models.ErrorResponse{
				Success: false,
				Error:   "User does not exist",
				Code:    http.StatusNotFound,
			})
			return
		}
		// Handle other database errors
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Database error",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Verify Password
	if err := utils.VerifyPassword(inputUser.Password, user.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Invalid password",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	// Generate Token
	token, err := middleware.GenerateToken(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to generate token",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	var userEmail string = user.Email

	ctx.JSON(http.StatusOK, models.LoginPhoneOutput{
		Phone: user.Phone,
		Email: userEmail,
		Token: token,
	})
}
