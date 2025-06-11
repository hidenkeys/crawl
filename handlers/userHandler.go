package handlers

import (
	"crawl/api"
	"crawl/models"
	"crawl/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"golang.org/x/crypto/bcrypt"
)

func (s Server) GetArtistDashboard(c *fiber.Ctx, artistId openapi_types.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetArtistRevenueDashboard(c *fiber.Ctx, artistId openapi_types.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) UserLogin(c *fiber.Ctx) error {
	var loginReq api.UserLoginJSONRequestBody
	if err := c.BodyParser(&loginReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Attempt to get the user by either username or email
	user, err := s.usrService.GetUserByUsername(c.Context(), loginReq.Username)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Compare the password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating token"})
	}

	// Return the token and user data
	return c.JSON(fiber.Map{
		"token": token,
		"user":  user,
	})
}

func (s Server) CreateUser(c *fiber.Ctx) error {
	var reqBody api.CreateUserJSONRequestBody
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ErrorCode": "400",
			"Message":   "Invalid request body",
		})
	}

	// Check if password and confirm_password match
	if reqBody.Password != reqBody.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ErrorCode": "400",
			"Message":   "Password and Confirm Password do not match",
		})
	}

	// Check if username already exists
	if existingUser, _ := s.usrService.GetUserByUsername(c.Context(), reqBody.Username); existingUser != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ErrorCode": "400",
			"Message":   "Username already exists",
		})
	}

	// Check if email already exists
	if existingUser, _ := s.usrService.GetUserByEmail(c.Context(), reqBody.Email); existingUser != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ErrorCode": "400",
			"Message":   "Email already exists",
		})
	}

	// Hash the password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error hashing password",
		})
	}

	// Create the user model
	user := &models.User{
		FirstName: reqBody.FirstName,
		LastName:  reqBody.LastName,
		Username:  reqBody.Username,
		Email:     reqBody.Email,
		Password:  string(hashedPassword),
		Role:      "listener", // Default role; adjust if needed
	}

	// Create the user using the service
	createdUser, err := s.usrService.CreateUser(c.Context(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ErrorCode": "500",
			"Message":   "Error creating user",
		})
	}

	// Return the created user info (omit password in response ideally)
	createdUser.Password = "" // Clear password before returning

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"Message": "User created successfully",
		"Data":    createdUser,
	})
}

func (s Server) DeleteUser(c *fiber.Ctx, userId openapi_types.UUID) error {
	// Convert the userId to uuid
	userUUID, err := uuid.Parse(userId.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Call the service to delete the user
	err = s.usrService.DeleteUser(c.Context(), userUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Return a success response
	return c.SendStatus(fiber.StatusNoContent)
}

func (s Server) GetUser(c *fiber.Ctx, userId openapi_types.UUID) error {
	// Convert userId to uuid
	userUUID, err := uuid.Parse(userId.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Call the service to fetch user details
	user, err := s.usrService.GetUserByID(c.Context(), userUUID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Return the user data
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Message": "Success",
		"Data":    user,
	})
}

func (s Server) UpdateUser(c *fiber.Ctx, userId openapi_types.UUID) error {
	var updateUserReq api.UpdateUserJSONRequestBody
	if err := c.BodyParser(&updateUserReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userUUID, err := uuid.Parse(userId.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Fetch the existing user from DB
	existingUser, err := s.usrService.GetUserByID(c.Context(), userUUID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Update only fields that are non-empty in request
	if updateUserReq.FirstName != "" {
		existingUser.FirstName = updateUserReq.FirstName
	}
	if updateUserReq.LastName != "" {
		existingUser.LastName = updateUserReq.LastName
	}
	if updateUserReq.Email != "" {
		existingUser.Email = updateUserReq.Email
	}
	if updateUserReq.Password != "" {
		// Hash the password before updating
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateUserReq.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
		}
		existingUser.Password = string(hashedPassword)
	}

	// Now save the updated user
	updatedUser, err := s.usrService.UpdateUser(c.Context(), userUUID, existingUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(updatedUser)
}
