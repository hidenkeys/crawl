package services

import (
	"context"
	"crawl/api"
	"crawl/models"
	"crawl/repositories"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"github.com/oapi-codegen/runtime/types"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type AuthService interface {
	Login(ctx context.Context, credentials api.PostLoginJSONBody) (*AuthResponse, error)
	ParseToken(tokenString string) (*Claims, error)
	generateJWTToken(user *models.User) (string, error)
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

type Claims struct {
	UserID types.UUID `json:"user_id"`
	Email  string     `json:"email"`
	Roles  []string   `json:"role"`
	jwt.RegisteredClaims
}

type authService struct {
	userRepo    repositories.IUserRepository
	jwtSecret   string
	tokenExpiry time.Duration
}

func NewAuthService(userRepo repositories.IUserRepository) AuthService {
	// Read JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET environment variable not set")
	}

	// Set token expiry (default to 24 hours)
	tokenExpiry := 24 * time.Hour
	if expiryStr := os.Getenv("JWT_EXPIRY"); expiryStr != "" {
		if duration, err := time.ParseDuration(expiryStr); err == nil {
			tokenExpiry = duration
		}
	}

	return &authService{
		userRepo:    userRepo,
		jwtSecret:   jwtSecret,
		tokenExpiry: tokenExpiry,
	}
}

func (s *authService) Login(ctx context.Context, credentials api.PostLoginJSONBody) (*AuthResponse, error) {
	// 1. Find user by email
	user, err := s.userRepo.FindByEmail(string(credentials.Email))
	if err != nil {

		log.Warn("Invalid credentials")
		return nil, errors.New("invalid email or password")
	}

	// 2. Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(credentials.Password)); err != nil {
		log.Warnf("Failed to hash password: %s", err.Error())
		return nil, errors.New("invalid email or password")
	}

	// 3. Generate JWT token
	token, err := s.generateJWTToken(user)
	if err != nil {
		log.Warn("Failed to generate token")
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 4. Return response
	return &AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *authService) generateJWTToken(user *models.User) (string, error) {
	// Create the claims
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"exp":     time.Now().Add(s.tokenExpiry).Unix(),
		"iat":     time.Now().Unix(),
		"role":    user.Roles,
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate signed token
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		log.Warn("Failed to generate signed token")
		return "", err
	}

	return signedToken, nil
}

func (s *authService) ParseToken(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate the claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Check if token is expired
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, errors.New("token is expired")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
