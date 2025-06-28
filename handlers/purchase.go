package handlers

import (
	"crawl/api"
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) PostPurchasesAlbums(c *fiber.Ctx) error {
	userID, err := h.getUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var purchaseReq api.PostPurchasesAlbumsJSONBody
	if err := c.BodyParser(&purchaseReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	// Verify the requesting user is purchasing for themselves
	if userID != purchaseReq.UserId {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can only purchase for yourself",
		})
	}

	// Verify the album exists
	_, err = h.Album.GetAlbumByID(c.Context(), purchaseReq.AlbumId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Album not found",
		})
	}

	// Process payment
	purchase, err := h.Purchase.PurchaseAlbum(c.Context(), purchaseReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to process purchase",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(purchase)
}

func (h *Handlers) PostPurchasesSongs(c *fiber.Ctx) error {
	userID, err := h.getUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var purchaseReq api.PostPurchasesSongsJSONBody
	if err := c.BodyParser(&purchaseReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	// Verify the requesting user is purchasing for themselves
	if userID != purchaseReq.UserId {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can only purchase for yourself",
		})
	}

	// Verify the song exists
	_, err = h.Song.GetSongByID(c.Context(), purchaseReq.SongId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Song not found",
		})
	}

	// Process payment
	purchase, err := h.Purchase.PurchaseSong(c.Context(), purchaseReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to process purchase",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(purchase)
}
