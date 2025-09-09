package utils

import (
	"errors"
	"net/http"
	"regexp"
	"tutuplapak/internal/models"
)

// Validator
func Validate(input *models.InputUser) error {

	// Email Validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(input.Email) {
		return errors.New("email format is invalid")
	}

	// Password Validation
	if len(input.Password) < 8 || len(input.Password) > 32 {
		return errors.New("password length must be 8–32 characters")
	}

	// Password Check Using Regex
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(input.Password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(input.Password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(input.Password)
	// hasSpecial := regexp.MustCompile(`[!@#$%^&*]`).MatchString(input.Password)

	// if !hasNumber || !hasUpper || !hasLower || !hasSpecial {
	// 	return errors.New("password must contain at least one number, uppercase letter, lowercase letter, and special character")
	// }
	// return nil

	if !hasNumber || !hasUpper || !hasLower {
		return errors.New("password must contain at least one number, uppercase letter, lowercase letter")
	}
	return nil
}

func PhoneValidation(phone string) *models.ErrorResponse {
	phoneRegex := regexp.MustCompile(`^\+[1-9]{1,3}[0-9]{7,14}$`)
	if !phoneRegex.MatchString(phone) {
		return &models.ErrorResponse{
			Success: false,
			Error:   "Invalid phone number",
			Code:    http.StatusBadRequest,
		}
	}
	return nil
}

func PasswordValidation(password string) *models.ErrorResponse {
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	// hasSpecial := regexp.MustCompile(`[!@#$%^&*]`).MatchString(password)

	if !hasNumber || !hasUpper || !hasLower {
		return &models.ErrorResponse{
			Success: false,
			Error:   "Password must contain at least one number, uppercase letter, lowercase letter",
			Code:    http.StatusBadRequest,
		}
	}
	return nil
}
