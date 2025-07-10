package handlers

import (
	"crawl/api"
	"crawl/models"
	"github.com/gofiber/fiber/v2"
	"github.com/oapi-codegen/runtime/types"
)

func (h *Handlers) GetAlbums(c *fiber.Ctx, params api.GetAlbumsParams) error {
	albums, err := h.Album.GetAllAlbums(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to fetch albums",
		})
	}

	return c.JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "Albums fetched successfully",
		Data:    albums,
	})
}

func (h *Handlers) PostAlbums(c *fiber.Ctx) error {
	// JWT authentication check
	userDetails, err := h.getDetailsFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var albumReq api.Album
	if err := c.BodyParser(&albumReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	// Verify the requesting user is the artist
	artist, err := h.User.GetArtistByUserId(c.Context(), userDetails.userID)
	if err != nil || artist.ID != albumReq.ArtistId {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can only create albums for yourself",
		})
	}
	modelAlbum := models.Album{
		Title:    albumReq.Title,
		ArtistID: albumReq.ArtistId,
	}

	// Handle nullable/optional fields
	if albumReq.Description != nil {
		modelAlbum.Description = *albumReq.Description
	}

	if albumReq.CoverImageUrl != nil {
		modelAlbum.CoverImageURL = *albumReq.CoverImageUrl
	}

	if albumReq.Price != nil {
		modelAlbum.Price = *albumReq.Price
	} else {
		modelAlbum.Price = 0 // Default value
	}

	if albumReq.ReleaseDate != nil {
		modelAlbum.ReleaseDate = *albumReq.ReleaseDate
	}

	if albumReq.IsFlagged != nil {
		modelAlbum.IsFlagged = *albumReq.IsFlagged
	} else {
		modelAlbum.IsFlagged = false // Default from gorm tag
	}

	// Handle ID if needed
	if albumReq.Id != nil {
		modelAlbum.ID = *albumReq.Id
	}

	createdAlbum, err := h.Album.CreateAlbum(c.Context(), modelAlbum)
	if err != nil {
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to create album",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "Album created successfully",
		Data:    createdAlbum,
	})
}

func (h *Handlers) GetAlbumsAlbumId(c *fiber.Ctx, albumId types.UUID) error {
	album, err := h.Album.GetAlbumByID(c.Context(), albumId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "Album not found",
		})
	}

	return c.JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "Albums fetched successfully",
		Data:    album,
	})
}

func (h *Handlers) PutAlbumsAlbumId(c *fiber.Ctx, albumId types.UUID) error {
	userDetails, err := h.getDetailsFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var albumReq api.Album
	if err := c.BodyParser(&albumReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	// Verify the requesting user is the album's artist
	album, err := h.Album.GetAlbumByID(c.Context(), albumId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "Album not found",
		})
	}

	artist, err := h.Artist.GetArtistByID(c.Context(), userDetails.userID)
	if err != nil || artist.ID != album.ArtistID {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can only update your own albums",
		})
	}

	modelAlbum := &models.Album{
		Title:    albumReq.Title,
		ArtistID: albumReq.ArtistId,
	}

	// Handle nullable/optional fields
	if albumReq.Description != nil {
		modelAlbum.Description = *albumReq.Description
	}

	if albumReq.CoverImageUrl != nil {
		modelAlbum.CoverImageURL = *albumReq.CoverImageUrl
	}

	if albumReq.Price != nil {
		modelAlbum.Price = *albumReq.Price
	} else {
		modelAlbum.Price = 0 // Default value
	}

	if albumReq.ReleaseDate != nil {
		modelAlbum.ReleaseDate = *albumReq.ReleaseDate
	}

	if albumReq.IsFlagged != nil {
		modelAlbum.IsFlagged = *albumReq.IsFlagged
	} else {
		modelAlbum.IsFlagged = false // Default from gorm tag
	}

	// Handle ID if needed
	if albumReq.Id != nil {
		modelAlbum.ID = *albumReq.Id
	}

	updatedAlbum, err := h.Album.UpdateAlbum(c.Context(), albumId, modelAlbum)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update album",
		})
	}

	return c.JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "Album updated successfully",
		Data:    updatedAlbum,
	})
}

func (h *Handlers) DeleteAlbumsAlbumId(c *fiber.Ctx, albumId types.UUID) error {
	userDetails, err := h.getDetailsFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	// Verify the requesting user is the album's artist
	album, err := h.Album.GetAlbumByID(c.Context(), albumId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "Album not found",
		})
	}

	artist, err := h.User.GetArtistByUserId(c.Context(), userDetails.userID)
	if err != nil || artist.ID != album.ArtistID {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can only delete your own albums",
		})
	}

	if err := h.Album.DeleteAlbum(c.Context(), albumId); err != nil {
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to delete album",
		})
	}

	return c.JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "Album deleted successfully",
		Data:    nil,
	})
}

func (h *Handlers) GetAlbumsAlbumIdSongs(c *fiber.Ctx, albumId types.UUID) error {
	songs, err := h.Album.GetAlbumSongs(c.Context(), albumId)
	if err != nil {
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to fetch album songs",
		})
	}

	return c.JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "Album songs fetched successfully",
		Data:    songs,
	})
}

func (h *Handlers) GetAlbumsAlbumIdContributors(c *fiber.Ctx, albumId types.UUID) error {
	contributors, err := h.Album.GetAlbumContributors(c.Context(), albumId)
	if err != nil {
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to fetch album contributors",
		})
	}

	return c.JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "Album contributors fetched successfully",
		Data:    contributors,
	})
}

func (h *Handlers) PostAlbumsAlbumIdContributors(c *fiber.Ctx, albumId types.UUID) error {
	userDetails, err := h.getDetailsFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var contributorReq api.Contributor
	if err := c.BodyParser(&contributorReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	// Verify the requesting user is the album's artist
	album, err := h.Album.GetAlbumByID(c.Context(), albumId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "Album not found",
		})
	}

	artist, err := h.User.GetArtistByUserId(c.Context(), userDetails.userID)
	if err != nil || artist.ID != album.ArtistID {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can only add contributors to your own albums",
		})
	}

	// Verify the contributor artist exists
	_, err = h.Artist.GetArtistByID(c.Context(), contributorReq.ArtistId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Contributor artist not found",
		})
	}

	contributor := &models.AlbumContributor{
		AlbumID:          album.ID,
		ArtistID:         contributorReq.ArtistId,
		ContributionType: contributorReq.ContributionType,
	}

	if contributorReq.RoyaltyPercentage != nil {
		contributor.RoyaltyPercentage = *contributorReq.RoyaltyPercentage
	}

	err = h.Album.AddAlbumContributor(c.Context(), albumId, contributor)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to add contributor",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "Album Contributor added successfully",
		Data:    nil,
	})
}

func (h *Handlers) PostAlbumsLikeAlbumId(c *fiber.Ctx, albumId api.AlbumId) error {
	_, err := h.getDetailsFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}
	err = h.Album.LikeAlbum(c.Context(), albumId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "Album not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "Album liked successfully",
		Data:    nil,
	})
}

func (h *Handlers) PostAlbumsBulk(c *fiber.Ctx) error {
	// JWT authentication check
	userDetails, err := h.getDetailsFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var bulkReq api.BulkAlbumRequest
	if err := c.BodyParser(&bulkReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	// Verify the requesting user is the artist for the album
	artist, err := h.User.GetArtistByUserId(c.Context(), userDetails.userID)
	if err != nil || artist.ID != bulkReq.Album.ArtistId {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can only create albums for yourself",
		})
	}

	// Convert api.Album to models.Album
	modelAlbum := models.Album{
		Title:    bulkReq.Album.Title,
		ArtistID: bulkReq.Album.ArtistId,
	}
	if bulkReq.Album.Description != nil {
		modelAlbum.Description = *bulkReq.Album.Description
	}
	if bulkReq.Album.CoverImageUrl != nil {
		modelAlbum.CoverImageURL = *bulkReq.Album.CoverImageUrl
	}
	if bulkReq.Album.Price != nil {
		modelAlbum.Price = *bulkReq.Album.Price
	} else {
		modelAlbum.Price = 0
	}
	if bulkReq.Album.ReleaseDate != nil {
		modelAlbum.ReleaseDate = *bulkReq.Album.ReleaseDate
	}
	if bulkReq.Album.IsFlagged != nil {
		modelAlbum.IsFlagged = *bulkReq.Album.IsFlagged
	} else {
		modelAlbum.IsFlagged = false
	}
	if bulkReq.Album.Id != nil {
		modelAlbum.ID = *bulkReq.Album.Id
	}

	// Convert api.Song slice to models.Song slice
	modelSongs := make([]models.Song, len(bulkReq.Songs))
	for i, s := range bulkReq.Songs {
		modelSongs[i] = models.Song{
			Title:       s.Title,
			ArtistID:    s.ArtistId,
			AlbumID:     s.AlbumId,
			Price:       s.Price,
			AudioURL:    s.AudioUrl,
			ReleaseDate: s.ReleaseDate,
			GenreID:     &s.GenreId,
			Lyrics:      s.Lyrics,
			Order:       s.Order,
		}
		if s.PreviewUrl != nil {
			modelSongs[i].PreviewURL = *s.PreviewUrl
		}
		if s.CoverImageUrl != nil {
			modelSongs[i].CoverImageURL = *s.CoverImageUrl
		}
	}

	createdAlbum, err := h.Album.CreateAlbumWithSongs(c.Context(), &modelAlbum, modelSongs)
	if err != nil {
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to create album and songs: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "Album and songs created successfully",
		Data:    createdAlbum,
	})
}
