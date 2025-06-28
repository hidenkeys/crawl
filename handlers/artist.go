package handlers

import (
	"crawl/api"
	"crawl/models"
	"github.com/gofiber/fiber/v2"
	"github.com/oapi-codegen/runtime/types"
)

func (h *Handlers) GetArtists(c *fiber.Ctx, params api.GetArtistsParams) error {
	artists, err := h.Artist.GetAllArtists(c.Context(), params.Page, params.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to fetch artists",
		})
	}

	return c.JSON(artists)
}

func (h *Handlers) PostArtists(c *fiber.Ctx) error {
	userID, err := h.getUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var artistReq api.Artist
	if err := c.BodyParser(&artistReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	// Check if user already has an artist profile

	_, err = h.User.GetArtistByUserId(c.Context(), userID)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "User already has an artist profile",
		})
	}

	artist := &models.Artist{
		ArtistName: artistReq.ArtistName,
		UserID:     userID,
	}

	if artistReq.Verified != nil {
		artist.Verified = *artistReq.Verified
	}

	if artistReq.WalletBalance != nil {
		artist.WalletBalance = float64(*artistReq.WalletBalance)
	}

	if artistReq.MonthlyListeners != nil {
		artist.MonthlyListeners = *artistReq.MonthlyListeners
	}

	if artistReq.Id != nil {
		artist.ID = *artistReq.Id
	}

	createdArtist, err := h.Artist.CreateArtist(c.Context(), artist)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create artist profile",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdArtist)
}

func (h *Handlers) GetArtistsArtistId(c *fiber.Ctx, artistId types.UUID) error {
	artist, err := h.Artist.GetArtistByID(c.Context(), artistId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "Artist not found",
		})
	}

	return c.JSON(artist)
}

func (h *Handlers) PutArtistsArtistId(c *fiber.Ctx, artistId types.UUID) error {
	userID, err := h.getUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var artistReq api.Artist
	if err := c.BodyParser(&artistReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	_, err = h.Artist.GetArtistByID(c.Context(), artistId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "Artist not Found",
		})
	}

	// Verify the requesting user owns the artist profile
	artist, err := h.User.GetArtistByUserId(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "Artist not an Artist",
		})
	}
	artist.ArtistName = artistReq.ArtistName

	if artistReq.Verified != nil {
		artist.Verified = *artistReq.Verified
	}

	if artistReq.WalletBalance != nil {
		artist.WalletBalance = float64(*artistReq.WalletBalance)
	}

	if artistReq.MonthlyListeners != nil {
		artist.MonthlyListeners = *artistReq.MonthlyListeners
	}

	if artistReq.Id != nil {
		artist.ID = *artistReq.Id
	}

	updatedArtist, err := h.Artist.UpdateArtist(c.Context(), artist.ID, artist)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update artist",
		})
	}

	return c.JSON(updatedArtist)
}

func (h *Handlers) GetArtistsArtistIdSongs(c *fiber.Ctx, artistId types.UUID, params api.GetArtistsArtistIdSongsParams) error {
	songs, err := h.Artist.GetArtistSongs(c.Context(), artistId, params.Page, params.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to fetch artist songs",
		})
	}

	return c.JSON(songs)
}
