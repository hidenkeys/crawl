package handlers

import (
	"crawl/api"
	"crawl/models"
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) PostStreams(c *fiber.Ctx) error {
	userDetails, err := h.getDetailsFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var streamReq api.PostStreamsJSONBody
	if err := c.BodyParser(&streamReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	// Verify the song exists
	_, err = h.Song.GetSongByID(c.Context(), streamReq.SongId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Song not found",
		})
	}

	// Record the stream
	stream := models.Stream{
		SongID:    streamReq.SongId,
		IsPreview: false, // default
	}

	*stream.UserID = userDetails.userID

	if streamReq.CountryCode != nil {
		stream.CountryCode = *streamReq.CountryCode
	}

	if streamReq.DeviceType != nil {
		stream.DeviceType = *streamReq.DeviceType
	}

	if streamReq.IsPreview != nil {
		stream.IsPreview = *streamReq.IsPreview
	}

	if err := h.Stream.RecordStream(c.Context(), stream); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to record stream",
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}
