package middleware

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	"tutuplapak/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

// Generate Token for Login and Register
func GenerateToken(user *models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// log.Println("JWT_SECRET not set in environment")
		return "", fmt.Errorf("JWT_SECRET not set")
	}

	expTime := time.Now().Add(30 * time.Minute)

	claims := &models.JWTClaim{
		ID:    user.ID,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   string(rune(user.ID)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Println("Failed to sign JWT:", err)
		return "", err
	}

	return signed, nil
}

// Parse Token
func ParseToken(tokenString string) (*models.JWTClaim, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET not set")
	}

	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.JWTClaim)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
