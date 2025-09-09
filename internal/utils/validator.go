package utils

import (
	"errors"
	"regexp"
	"tutuplapak/internal/models"
)

func EmailValidation(emailInput string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(emailInput) {
		return errors.New("email format is invalid")
	}

	return nil
}

func PasswordLengthValidation(password string) error {
	// Password Validation
	if len(password) < 8 || len(password) > 32 {
		return errors.New("password length must be 8–32 characters")
	}
	return nil
}

// Validator
func Validate(input *models.InputUser) error {

	// Email Validation

	err := EmailValidation(input.Email)
	if err != nil {
		return err
	}

	// Password Validation
	err = PasswordLengthValidation(input.Password)
	if err != nil {
		return err
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
