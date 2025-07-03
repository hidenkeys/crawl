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
	detailsFromToken, err := h.getDetailsFromToken(c)
	isAdmin := false
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
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
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}
	users, err := h.User.GetAllUsers(c.Context(), *params.Page, *params.Limit)
	if err != nil {
		log.Errorf("Failed to fetch users: %s", err.Error())
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "Failed to fetch users",
		})
	}

	return c.JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "Users Fetched successfully",
		Data:    users,
	})
}

func (h *Handlers) PostCreateAdmin(c *fiber.Ctx) error {
	detailsFromToken, err := h.getDetailsFromToken(c)
	isAdmin := false
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
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
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var userReq api.PostCreateAdminJSONBody
	if err := c.BodyParser(&userReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}
	err = h.User.CreateAdminFromUser(c.Context(), *userReq.UserId)
	if err != nil {
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to assign admin role to user",
		})
	}
	return c.Status(fiber.StatusOK).JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "Admin created successfully",
		Data:    nil,
	})

}

func (h *Handlers) PostUsers(c *fiber.Ctx) error {
	var userReq api.PostUsersJSONRequestBody
	if err := c.BodyParser(&userReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	// Check if email already exists
	_, err := h.User.GetUserByEmail(c.Context(), string(userReq.Email))
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(api.Error{
			Code:    fiber.StatusConflict,
			Message: "Email already exists",
		})
	}

	// Check if username already exists
	_, err = h.User.GetUserByUsername(c.Context(), userReq.Username)
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(api.Error{
			Code:    fiber.StatusConflict,
			Message: "Username already exists",
		})
	}
	println(userReq.LastName)

	pass, err := bcrypt.GenerateFromPassword([]byte(*userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("Failed to hash password: %s", err.Error())
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{ // 500
			Code:    fiber.StatusExpectationFailed,
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

	_, err = h.User.Create(c.Context(), newUser)
	if err != nil {
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to create user: " + err.Error(),
		})
	}
	var loginReq api.PostLoginJSONBody
	loginReq.Password = *userReq.Password
	loginReq.Email = userReq.Email

	token, err := h.Auth.Login(c.Context(), loginReq)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Invalid credentials",
		})
	}

	/*return c.JSON(fiber.Map{
		"token": token,
	})*/
	return c.Status(fiber.StatusOK).JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "Login Successful",
		Data:    token,
	})

	//return c.Status(fiber.StatusCreated).JSON(createdUser)
}

func (h *Handlers) GetUsersUserId(c *fiber.Ctx, userId types.UUID) error {
	detailsFromToken, err := h.getDetailsFromToken(c)
	isAdmin := false
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}
	for _, v := range detailsFromToken.roles {
		if v.Name == "Admin" {
			isAdmin = true
			_ = v
		}
	}

	// Verify the requesting user is updating their own profile
	if !isAdmin || detailsFromToken.userID != userId {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "Unauthorized",
		})
	}
	user, err := h.User.GetByID(c.Context(), userId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "User not found",
		})
	}

	return c.JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "User fetched successfully",
		Data:    user,
	})
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
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to update user",
		})
	}

	return c.JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "User updated successfully",
		Data:    updatedUser,
	})
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
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to delete user",
		})
	}

	return c.JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "User deleted successfully",
		Data:    nil,
	})
}

func (h *Handlers) GetUsersUserIdPlaylists(c *fiber.Ctx, userId types.UUID) error {
	requestingUserDetails, err := h.getDetailsFromToken(c)
	if err != nil {
		// If not authenticated, only show public playlists
		playlists, err := h.Playlist.GetAllPlaylists(c.Context())
		if err != nil {
			return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
				Code:    fiber.StatusExpectationFailed,
				Message: "Failed to fetch user playlists",
			})
		}
		return c.JSON(models.Response{
			Code:    fiber.StatusOK,
			Message: "Users playlists fetched successfully",
			Data:    playlists,
		})
	}

	// If authenticated and requesting own playlists, show all
	if requestingUserDetails.userID == userId {
		playlists, err := h.User.GetUserPlaylists(c.Context(), userId)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(api.Error{
				Code:    fiber.StatusNotFound,
				Message: "Failed to fetch user playlists",
			})
		}
		return c.JSON(models.Response{
			Code:    fiber.StatusOK,
			Message: "Playlists fetched successfully",
			Data:    playlists,
		})
	}

	// If authenticated but requesting someone else's playlists, show only public ones
	playlists, err := h.User.GetUserPublicPlaylists(c.Context(), userId)
	if err != nil {
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to fetch user playlists",
		})
	}
	return c.JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "Playlists fetched successfully",
		Data:    playlists,
	})
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
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to create playlist",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "Playlist created successfully",
		Data:    createdPlaylist,
	})
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
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to fetch purchased albums",
		})
	}

	return c.JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "Library albums fetched successfully",
		Data:    albums,
	})
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
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to fetch purchased songs",
		})
	}

	return c.JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "Library songs fetched successfully",
		Data:    songs,
	})
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
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to fetch purchase history",
		})
	}

	return c.JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "User purchases fetched successfully",
		Data:    purchases,
	})
}

func (h *Handlers) GetUsersUserIdRecent(c *fiber.Ctx, userId api.UserId, params api.GetUsersUserIdRecentParams) error {
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
			Message: "Forbidden, You can only update your own profile",
		})
	}

	songs, err := h.User.GetUserRecentStreams(c.Context(), userId, *params.Limit, *params.Page)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "An error occurred, songs not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "Users recently played songs fetched successfully",
		Data:    songs,
	})

}
