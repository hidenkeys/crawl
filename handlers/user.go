package handlers

import (
	"crawl/api"
	"crawl/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/oapi-codegen/runtime/types"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handlers) GetUsers(c *fiber.Ctx, params api.GetUsersParams) error {
	users, err := h.User.GetAllUsers(c.Context(), *params.Page, *params.Limit)
	if err != nil {
		log.Errorf("Failed to fetch users: %s", err.Error())
		return c.Status(fiber.StatusOK).JSON(api.Error{
			Code:    fiber.StatusOK,
			Message: "Failed to fetch users",
		})
	}

	return c.JSON(users)
}

func (h *Handlers) PostCreateAdmin(c *fiber.Ctx) error {
	detailsFromToken, err := h.getDetailsFromToken(c)
	isAdmin := false
	var _ models.Role
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(api.Error{
			Code:    fiber.StatusOK,
			Message: "Unauthorized",
		})
	}
	for _, v := range detailsFromToken.roles {
		if v.Name == "Admin" {
			isAdmin = true
			_ = v
		}
	}
	if !isAdmin {
		return c.Status(fiber.StatusOK).JSON(api.Error{
			Code:    fiber.StatusOK,
			Message: "Unauthorized",
		})
	}

	var userReq api.PostCreateAdminJSONBody
	if err := c.BodyParser(&userReq); err != nil {
		return c.Status(fiber.StatusOK).JSON(api.Error{
			Code:    01,
			Message: "Invalid request body",
		})
	}
	// TODO complete
	return c.Status(fiber.StatusOK).JSON("Done")

}

func (h *Handlers) PostUsers(c *fiber.Ctx) error {
	var userReq api.PostUsersJSONRequestBody
	if err := c.BodyParser(&userReq); err != nil {
		return c.Status(fiber.StatusOK).JSON(api.Error{
			Code:    01,
			Message: "Invalid request body",
		})
	}

	// Check if email already exists
	_, err := h.User.GetUserByEmail(c.Context(), string(userReq.Email))
	if err == nil {
		return c.Status(fiber.StatusOK).JSON(api.Error{
			Code:    02,
			Message: "Email already exists",
		})
	}

	// Check if username already exists
	_, err = h.User.GetUserByUsername(c.Context(), userReq.Username)
	if err == nil {
		return c.Status(fiber.StatusOK).JSON(api.Error{
			Code:    02,
			Message: "Username already exists",
		})
	}
	println(userReq.LastName)

	pass, err := bcrypt.GenerateFromPassword([]byte(*userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("Failed to hash password: %s", err.Error())
		return c.Status(fiber.StatusOK).JSON(api.Error{ // 500
			Code:    03,
			Message: "Failed to process password",
		})
	}
	newUser := &models.User{
		Email:          string(userReq.Email),
		HashedPassword: string(pass),
		Username:       userReq.Username,
		FirstName:      userReq.FirstName,
		LastName:       userReq.LastName,
	}

	// 6. Handle optional fields safely
	if userReq.Bio != nil {
		newUser.Bio = *userReq.Bio
	}
	if userReq.PhoneNumber != nil {
		newUser.PhoneNumber = *userReq.PhoneNumber
	}
	if userReq.ProfileImageUrl != nil {
		newUser.ProfileImage = *userReq.ProfileImageUrl
	}
	//newUser.Roles = userReq.
	//newUser.ArtistProfile

	println("After parse: " + newUser.LastName)

	createdUser, err := h.User.Create(c.Context(), newUser)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(api.Error{
			Code:    04,
			Message: "Failed to create user: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdUser)
}

func (h *Handlers) GetUsersUserId(c *fiber.Ctx, userId types.UUID) error {
	user, err := h.User.GetByID(c.Context(), userId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "User not found",
		})
	}

	return c.JSON(user)
}

func (h *Handlers) PutUsersUserId(c *fiber.Ctx, userId types.UUID) error {
	detailsFromToken, err := h.getDetailsFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	// Verify the requesting user is updating their own profile
	if detailsFromToken.userID != userId {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can only update your own profile",
		})
	}

	var userReq api.User
	if err := c.BodyParser(&userReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	updatedUser, err := h.User.Update(c.Context(), userId, userReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update user",
		})
	}

	return c.JSON(updatedUser)
}

func (h *Handlers) DeleteUsersUserId(c *fiber.Ctx, userId types.UUID) error {
	requestingUserDetails, err := h.getDetailsFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	// Verify the requesting user is deleting their own profile
	if requestingUserDetails.userID != userId {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can only delete your own profile",
		})
	}

	if err := h.User.Delete(c.Context(), userId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to delete user",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handlers) GetUsersUserIdPlaylists(c *fiber.Ctx, userId types.UUID) error {
	requestingUserDetails, err := h.getDetailsFromToken(c)
	if err != nil {
		// If not authenticated, only show public playlists
		playlists, err := h.Playlist.GetAllPlaylists(c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
				Code:    fiber.StatusInternalServerError,
				Message: "Failed to fetch user playlists",
			})
		}
		return c.JSON(playlists)
	}

	// If authenticated and requesting own playlists, show all
	if requestingUserDetails.userID == userId {
		playlists, err := h.User.GetUserPlaylists(c.Context(), userId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
				Code:    fiber.StatusInternalServerError,
				Message: "Failed to fetch user playlists",
			})
		}
		return c.JSON(playlists)
	}

	// If authenticated but requesting someone else's playlists, show only public ones
	playlists, err := h.User.GetUserPublicPlaylists(c.Context(), userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to fetch user playlists",
		})
	}
	return c.JSON(playlists)
}

func (h *Handlers) PostUsersUserIdPlaylists(c *fiber.Ctx, userId types.UUID) error {
	requestingUserDetails, err := h.getDetailsFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	// Verify the requesting user is creating a playlist for themselves
	if requestingUserDetails.userID != userId {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can only create playlists for yourself",
		})
	}

	var playlistReq api.Playlist
	if err := c.BodyParser(&playlistReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	newPlaylist := &models.Playlist{
		Title:  playlistReq.Title,
		UserID: playlistReq.UserId,
	}

	// Handle nullable/optional fields
	if playlistReq.Description != nil {
		newPlaylist.Description = *playlistReq.Description
	}

	if playlistReq.CoverImageUrl != nil {
		newPlaylist.CoverImageURL = *playlistReq.CoverImageUrl
	}

	if playlistReq.IsPublic != nil {
		newPlaylist.IsPublic = *playlistReq.IsPublic
	} else {
		newPlaylist.IsPublic = false // Default value from gorm tag
	}

	if playlistReq.UpdatedAt != nil {
		newPlaylist.UpdatedAt = *playlistReq.UpdatedAt
	}

	// Handle ID if needed (assuming BaseModel has ID field)
	if playlistReq.Id != nil {
		newPlaylist.ID = *playlistReq.Id
	}

	createdPlaylist, err := h.User.CreatePlaylist(c.Context(), userId, newPlaylist)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create playlist",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdPlaylist)
}

func (h *Handlers) GetUsersUserIdLibraryAlbums(c *fiber.Ctx, userId types.UUID) error {
	requestingUserDetails, err := h.getDetailsFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	// Verify the requesting user is accessing their own library
	if requestingUserDetails.userID != userId {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can only access your own library",
		})
	}

	albums, err := h.User.GetUserPurchasedAlbums(c.Context(), userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to fetch purchased albums",
		})
	}

	return c.JSON(albums)
}

func (h *Handlers) GetUsersUserIdLibrarySongs(c *fiber.Ctx, userId types.UUID) error {
	requestingUserDetails, err := h.getDetailsFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	// Verify the requesting user is accessing their own library
	if requestingUserDetails.userID != userId {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can only access your own library",
		})
	}

	songs, err := h.User.GetUserPurchasedSongs(c.Context(), userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to fetch purchased songs",
		})
	}

	return c.JSON(songs)
}

func (h *Handlers) GetUsersUserIdLibraryPurchases(c *fiber.Ctx, userId types.UUID, params api.GetUsersUserIdLibraryPurchasesParams) error {
	requestingUserDetails, err := h.getDetailsFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	// Verify the requesting user is accessing their own library
	if requestingUserDetails.userID != userId {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can only access your own library",
		})
	}

	purchases, err := h.User.GetUserPurchaseHistory(c.Context(), userId, params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to fetch purchase history",
		})
	}

	return c.JSON(purchases)
}
