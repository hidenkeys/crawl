package handlers

import (
	"crawl/api"
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) PostTips(c *fiber.Ctx) error {
	userDetails, err := h.getDetailsFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var tipReq api.PostTipsJSONBody
	if err := c.BodyParser(&tipReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	// Verify the artist exists
	_, err = h.Artist.GetArtistByID(c.Context(), tipReq.ArtistId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Artist not found",
		})
	}

	// Process tip payment
	tip, err := h.Tip.SendTip(c.Context(), tipReq, userDetails.userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to process tip",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(tip)
}
