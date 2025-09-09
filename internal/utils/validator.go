package utils

import (
	"errors"
	"regexp"
	"strings"
	"tutuplapak/internal/models"
)

var e164Regex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

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
func Validate(input *models.LoginEmailInput) error {

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

func PhoneValidation(input *models.PhoneUser) error {
	phone := strings.TrimSpace(input.Phone)

	if phone == "" {
		return errors.New("phone number is required")
	}

	if !e164Regex.MatchString(phone) {
		return errors.New("invalid phone number format")
	}

	return nil
}

func PasswordValidation(input *models.PhoneUser) error {
	password := input.Password

	if len(input.Password) < 8 || len(input.Password) > 32 {
		return errors.New("password length must be 8–32 characters")
	}

	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)

	if !hasNumber || !hasUpper || !hasLower {
		return errors.New("password must contain at least one number, one uppercase letter, and one lowercase letter")
	}

	return nil
}
