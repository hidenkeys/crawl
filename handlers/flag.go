package handlers

import (
	"crawl/api"
	"crawl/models"
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) PostFlags(c *fiber.Ctx) error {
	userID, err := h.getUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var flagReq api.PostFlagsJSONBody
	if err := c.BodyParser(&flagReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	// Verify the target exists
	switch flagReq.TargetType {
	case api.PostFlagsJSONBodyTargetTypeAlbum:
		_, err = h.Album.GetAlbumByID(c.Context(), flagReq.TargetId)
	case api.PostFlagsJSONBodyTargetTypeSong:
		_, err = h.Song.GetSongByID(c.Context(), flagReq.TargetId)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid target type",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Target not found",
		})
	}

	// Create the flag
	flag := &models.ContentFlag{
		ReporterUserID: userID,
		TargetID:       flagReq.TargetId,
		TargetType:     string(flagReq.TargetType),
		Reason:         flagReq.Reason,
		Status:         "pending",
	}

	if flagReq.Description != nil {
		flag.Description = *flagReq.Description
	}

	newFlag, err := h.Moderation.FlagContent(c.Context(), flag)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to create flag",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(newFlag)
}
