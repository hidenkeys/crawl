package handlers

import (
	"crawl/api"
	"crawl/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"strconv"
	"strings"
)

// GetPurchasedSongs returns all purchased songs by a user
func (s Server) GetPurchasedSongs(c *fiber.Ctx, userId openapi_types.UUID, params api.GetPurchasedSongsParams) error {
	userUUID, err := uuid.Parse(userId.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ErrorCode": "400",
			"Message":   "Invalid user ID",
		})
	}

	limit := 10
	offset := 0

	if params.Limit != nil {
		limit = *params.Limit
	}
	if params.Offset != nil {
		offset = *params.Offset
	}

	songs, count, err := s.purchaseService.GetPurchasesByUser(c.Context(), userUUID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ErrorCode": "500",
			"Message":   "Failed to fetch purchased songs",
		})
	}

	return c.JSON(fiber.Map{
		"Data":    songs,
		"Message": "success",
		"Count":   count,
	})
}

// GetAllSongs returns a paginated list of all songs
func (s Server) GetAllSongs(c *fiber.Ctx, params api.GetAllSongsParams) error {
	limitStr := c.Query("limit", "20")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	//songs, total, err := s.songService.(c.Context(), limit, offset)
	//if err != nil {
	//	return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
	//		"ErrorCode": "500",
	//		"Message":   "Failed to fetch songs",
	//	})
	//}

	// Return with pagination info
	return c.JSON(fiber.Map{
		//"total":  total,
		"limit":  limit,
		"offset": offset,
		//"data":   songs,
	})
}

// CreateSong creates a new song record
func (s Server) CreateSong(c *fiber.Ctx) error {
	var req api.CreateSongJSONRequestBody
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ErrorCode": "400",
			"Message":   "Invalid request body",
		})
	}

	fmt.Println("in handler CreateSong")
	fmt.Print(req.ArtistsNames)

	// Map API request to model
	song := &models.Song{
		ID:           uuid.New(),
		Title:        req.Title,
		ArtistID:     uuid.MustParse(req.ArtistId.String()),
		ArtistsNames: pq.StringArray(req.ArtistsNames),
		Genre:        req.Genre,
		Price:        100,
		Duration:     *req.Duration,
		AudioURL:     *req.AudioUrl,
		ReleaseDate:  *req.ReleaseDate,
	}

	if req.AlbumId != nil {
		song.AlbumID = uuid.MustParse(req.AlbumId.String())
	}

	createdSong, err := s.songService.CreateSong(c.Context(), song)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ErrorCode": "500",
			"Message":   "Failed to create song",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdSong)
}

// SearchSongs searches songs by title, artist, or genre based on query param
func (s Server) SearchSongs(c *fiber.Ctx, params api.SearchSongsParams) error {
	query := strings.TrimSpace(params.Query)
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ErrorCode": "400",
			"Message":   "Query parameter is required",
		})
	}
	//
	//limit := 20
	//offset := 0
	//if params.Limit != nil && *params.Limit > 0 {
	//	limit = int(*params.Limit)
	//}
	//if params.Offset != nil && *params.Offset >= 0 {
	//	offset = int(*params.Offset)
	//}

	//results, total, err := s.songService.SearchSongs(c.Context(), query, limit, offset)
	//if err != nil {
	//	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	//		"ErrorCode": "500",
	//		"Message":   "Failed to search songs",
	//	})
	//}
	//
	//return c.JSON(fiber.Map{
	//	"total":  total,
	//	"limit":  limit,
	//	"offset": offset,
	//	"data":   results,
	//})
	return nil
}

// DeleteSong deletes a song by ID
func (s Server) DeleteSong(c *fiber.Ctx, songId openapi_types.UUID) error {
	songUUID, err := uuid.Parse(songId.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ErrorCode": "400",
			"Message":   "Invalid song ID",
		})
	}

	err = s.songService.DeleteSong(c.Context(), songUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ErrorCode": "500",
			"Message":   "Failed to delete song",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetSong retrieves a song by ID
func (s Server) GetSong(c *fiber.Ctx, songId openapi_types.UUID) error {
	songUUID, err := uuid.Parse(songId.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ErrorCode": "400",
			"Message":   "Invalid song ID",
		})
	}

	song, err := s.songService.GetSongByID(c.Context(), songUUID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"ErrorCode": "404",
			"Message":   "Song not found",
		})
	}

	return c.JSON(song)
}

// GetArtistSongs gets all songs by a specific artist
func (s Server) GetArtistSongs(c *fiber.Ctx, artistId openapi_types.UUID, params api.GetArtistSongsParams) error {
	artistUUID, err := uuid.Parse(artistId.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ErrorCode": "400",
			"Message":   "Invalid artist ID",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	songs, total, err := s.songService.GetSongsByArtist(c.Context(), artistUUID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ErrorCode": "500",
			"Message":   "Failed to fetch songs for artist",
		})
	}

	return c.JSON(fiber.Map{
		"total":  total,
		"limit":  limit,
		"offset": offset,
		"data":   songs,
	})
}

// UpdateSong updates song fields partially, hashing password if updated
func (s Server) UpdateSong(c *fiber.Ctx, songId openapi_types.UUID) error {
	var updateReq api.UpdateSongJSONRequestBody
	if err := c.BodyParser(&updateReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ErrorCode": "400",
			"Message":   "Invalid request body",
		})
	}

	songUUID, err := uuid.Parse(songId.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ErrorCode": "400",
			"Message":   "Invalid song ID",
		})
	}

	// Get existing song
	existingSong, err := s.songService.GetSongByID(c.Context(), songUUID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"ErrorCode": "404",
			"Message":   "Song not found",
		})
	}

	// Update only provided fields
	if updateReq.Title != "" {
		existingSong.Title = updateReq.Title
	}

	if updateReq.ArtistId != nil {
		existingSong.ArtistID = uuid.MustParse(updateReq.ArtistId.String())
	}

	if updateReq.ArtistsNames != nil {
		existingSong.ArtistsNames = updateReq.ArtistsNames
	}

	if updateReq.Genre != "" {
		existingSong.Genre = updateReq.Genre
	}

	if updateReq.Price != nil {
		existingSong.Price = *updateReq.Price
	}

	if updateReq.Duration != nil {
		existingSong.Duration = *updateReq.Duration // Dereference the pointer
	}

	if updateReq.AudioUrl != nil {
		existingSong.AudioURL = *updateReq.AudioUrl
	}

	if updateReq.ReleaseDate != nil {
		existingSong.ReleaseDate = *updateReq.ReleaseDate
	}

	if updateReq.AlbumId != nil {
		existingSong.AlbumID = uuid.MustParse(updateReq.AlbumId.String())
	}

	updatedSong, err := s.songService.UpdateSong(c.Context(), songUUID, existingSong)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ErrorCode": "500",
			"Message":   "Failed to update song",
		})
	}

	return c.JSON(updatedSong)
}
