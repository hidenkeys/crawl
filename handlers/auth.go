package handlers

import (
	"crawl/api"
	"github.com/gofiber/fiber/v2/log"
	"github.com/oapi-codegen/runtime/types"
	"strings"

	"github.com/gofiber/fiber/v2"
)

//// Auth middleware to verify JWT tokens
//func (h *Handlers) authMiddleware(c *fiber.Ctx) error {
//	tokenString := c.Get("Authorization")
//	if tokenString == "" {
//		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
//			"error": "Authorization header missing",
//		})
//	}
//
//	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, errors.New("unexpected signing method")
//		}
//		return []byte(h.jwtSecret), nil
//	})
//
//	if err != nil || !token.Valid {
//		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
//			"error": "Invalid token",
//		})
//	}
//
//	// Set user claims in context for later use
//	claims := token.Claims.(jwt.MapClaims)
//	c.Locals("userID", claims["user_id"])
//	c.Locals("isArtist", claims["is_artist"])
//
//	return c.Next()
//}
//
//// Admin middleware for admin-only endpoints
//func (h *Handlers) adminMiddleware(c *fiber.Ctx) error {
//	// Check if user is admin
//	// Example:
//	// userID := c.Locals("userID").(string)
//	// isAdmin, err := h.userService.IsAdmin(userID)
//	// if err != nil || !isAdmin {
//	//     return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
//	//         "error": "Admin access required",
//	//     })
//	// }
//	return c.Next()
//}
//
//// Helper function to get user ID from JWT token
//func (h *Handlers) getUserIDFromToken(c *fiber.Ctx) (openapi_types.UUID, error) {
//	// Get the token from the Authorization header
//	authHeader := c.Get("Authorization")
//	if authHeader == "" {
//		return openapi_types.UUID{}, fiber.NewError(fiber.StatusUnauthorized, "Missing authorization header")
//	}
//
//	// The token is typically in the format "Bearer <token>"
//	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
//	if tokenString == authHeader {
//		return openapi_types.UUID{}, fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization header format")
//	}
//
//	// Parse and validate the token
//	claims, err := h.Auth.ParseToken(tokenString)
//	if err != nil {
//		return openapi_types.UUID{}, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
//	}
//
//	// Convert the user ID from string to UUID
//	userID, err := (claims.UserID)
//	if err != nil {
//		return openapi_types.UUID{}, fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID in token")
//	}
//
//	return userID, nil
//}

func (h *Handlers) PostLogin(c *fiber.Ctx) error {
	var loginReq api.PostLoginJSONBody
	if err := c.BodyParser(&loginReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	token, err := h.Auth.Login(c.Context(), loginReq)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Invalid credentials",
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}

func (h *Handlers) getUserIDFromToken(c *fiber.Ctx) (types.UUID, error) {
	// Get the token from the Authorization header
	authHeader := c.Get("Authorization")
	log.Infof("Header: %s", authHeader)
	if authHeader == "" {
		return types.UUID{}, fiber.NewError(fiber.StatusUnauthorized, "Missing authorization header")
	}

	// The token is typically in the format "Bearer <token>"
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return types.UUID{}, fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization header format")
	}

	// Parse and validate the token
	claims, err := h.Auth.ParseToken(tokenString)
	if err != nil {
		return types.UUID{}, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	// Convert the user ID from string to UUID
	userID := claims.UserID

	return userID, nil
}
