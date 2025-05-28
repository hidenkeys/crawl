package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

var jwtSecret = []byte("your_secret_key") // Use a secure secret in production

// GenerateJWT generates a JWT token with user info
func GenerateJWT(userID uuid.UUID, username string, role string) (string, error) {
	claims := jwt.MapClaims{
		"id":       userID.String(),
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
