package handlers

import (
	"crawl/api"
	"crawl/models"
	"github.com/gofiber/fiber/v2"
	"github.com/oapi-codegen/runtime/types"
)

func (h *Handlers) GetSongs(c *fiber.Ctx, params api.GetSongsParams) error {
	songs, err := h.Song.GetAllSongs(c.Context(), params.Page, params.Limit, params.Genre, params.Artist, params.Album)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to fetch songs",
		})
	}

	return c.JSON(songs)
}

func (h *Handlers) PostSongs(c *fiber.Ctx) error {
	userDetails, err := h.getDetailsFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var songReq api.Song
	if err := c.BodyParser(&songReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	// Verify the requesting user is the song's artist
	artist, err := h.User.GetArtistByUserId(c.Context(), userDetails.userID)
	if err != nil || artist.ID != songReq.ArtistId {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can only create songs for yourself",
		})
	}

	// Verify the album exists if provided
	if songReq.AlbumId != nil {
		_, err := h.Album.GetAlbumByID(c.Context(), *songReq.AlbumId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(api.Error{
				Code:    fiber.StatusBadRequest,
				Message: "Album not found",
			})
		}
	}

	// Verify the genre exists
	_, err = h.Genre.GetGenreByID(c.Context(), songReq.GenreId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Genre not found",
		})
	}

	song := &models.Song{
		Title:       songReq.Title,
		ArtistID:    songReq.ArtistId,
		Duration:    songReq.Duration,
		Price:       songReq.Price,
		AudioURL:    songReq.AudioUrl,
		ReleaseDate: songReq.ReleaseDate,
		GenreID:     &songReq.GenreId,
	}

	if songReq.AlbumId != nil {
		song.AlbumID = songReq.AlbumId
	}

	if songReq.CoverImageUrl != nil {
		song.CoverImageURL = *songReq.CoverImageUrl
	}

	if songReq.PreviewUrl != nil {
		song.PreviewURL = *songReq.PreviewUrl
	}

	if songReq.IsFlagged != nil {
		song.IsFlagged = *songReq.IsFlagged
	}

	if songReq.PlaysCount != nil {
		song.PlaysCount = *songReq.PlaysCount
	}

	if songReq.CreatedAt != nil {
		song.CreatedAt = *songReq.CreatedAt
	}

	if songReq.UpdatedAt != nil {
		song.UpdatedAt = *songReq.UpdatedAt
	}

	if songReq.Id != nil {
		song.ID = *songReq.Id
	}

	createdSong, err := h.Song.CreateSong(c.Context(), song)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create song",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdSong)
}

func (h *Handlers) GetSongsSongId(c *fiber.Ctx, songId types.UUID) error {
	song, err := h.Song.GetSongByID(c.Context(), songId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "Song not found",
		})
	}

	return c.JSON(song)
}

func (h *Handlers) PutSongsSongId(c *fiber.Ctx, songId types.UUID) error {
	userDetails, err := h.getDetailsFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var songReq api.Song
	if err := c.BodyParser(&songReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	// Verify the requesting user is the song's artist
	song, err := h.Song.GetSongByID(c.Context(), songId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "Song not found",
		})
	}

	artist, err := h.User.GetArtistByUserId(c.Context(), userDetails.userID)
	if err != nil || artist.ID != song.ArtistID {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can only update your own songs",
		})
	}

	song.Title = songReq.Title
	song.ArtistID = songReq.ArtistId
	song.Duration = songReq.Duration
	song.Price = songReq.Price
	song.AudioURL = songReq.AudioUrl
	song.ReleaseDate = songReq.ReleaseDate
	song.GenreID = &songReq.GenreId

	if songReq.AlbumId != nil {
		song.AlbumID = songReq.AlbumId
	}

	if songReq.CoverImageUrl != nil {
		song.CoverImageURL = *songReq.CoverImageUrl
	}

	if songReq.PreviewUrl != nil {
		song.PreviewURL = *songReq.PreviewUrl
	}

	if songReq.IsFlagged != nil {
		song.IsFlagged = *songReq.IsFlagged
	}

	if songReq.PlaysCount != nil {
		song.PlaysCount = *songReq.PlaysCount
	}

	if songReq.CreatedAt != nil {
		song.CreatedAt = *songReq.CreatedAt
	}

	if songReq.UpdatedAt != nil {
		song.UpdatedAt = *songReq.UpdatedAt
	}

	if songReq.Id != nil {
		song.ID = *songReq.Id
	}

	updatedSong, err := h.Song.UpdateSong(c.Context(), songId, song)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update song",
		})
	}

	return c.JSON(updatedSong)
}

func (h *Handlers) DeleteSongsSongId(c *fiber.Ctx, songId types.UUID) error {
	userDetails, err := h.getDetailsFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	// Verify the requesting user is the song's artist
	song, err := h.Song.GetSongByID(c.Context(), songId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "Song not found",
		})
	}

	artist, err := h.User.GetArtistByUserId(c.Context(), userDetails.userID)
	if err != nil || artist.ID != song.ArtistID {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can only delete your own songs",
		})
	}

	if err := h.Song.DeleteSong(c.Context(), songId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to delete song",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handlers) GetSongsSongIdContributors(c *fiber.Ctx, songId types.UUID) error {
	contributors, err := h.Song.GetSongContributors(c.Context(), songId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to fetch song contributors",
		})
	}

	return c.JSON(contributors)
}

func (h *Handlers) PostSongsSongIdContributors(c *fiber.Ctx, songId types.UUID) error {
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

	// Verify the requesting user is the song's artist
	song, err := h.Song.GetSongByID(c.Context(), songId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "Song not found",
		})
	}

	artist, err := h.User.GetArtistByUserId(c.Context(), userDetails.userID)
	if err != nil || artist.ID != song.ArtistID {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can only add contributors to your own songs",
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

	contributor := &models.SongContributor{
		SongID:           songId,
		ArtistID:         contributorReq.ArtistId,
		ContributionType: contributorReq.ContributionType,
	}

	if contributorReq.RoyaltyPercentage != nil {
		contributor.RoyaltyPercentage = *contributorReq.RoyaltyPercentage
	}

	err = h.Song.AddSongContributor(c.Context(), songId, contributor)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to add contributor",
		})
	}

	return c.Status(fiber.StatusCreated).JSON("Added Song contributor")
}
