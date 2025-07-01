package handlers

import (
	"crawl/api"
	"crawl/models"
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) PostFlags(c *fiber.Ctx) error {
	userDetails, err := h.getDetailsFromToken(c)
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
		ReporterUserID: userDetails.userID,
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

func (h *Handlers) GetFlags(c *fiber.Ctx) error {
	detailsFromToken, err := h.getDetailsFromToken(c)
	isAdmin := false
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(api.Error{
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
		return c.Status(fiber.StatusOK).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	flags, err := h.Moderation.GetFlaggedContent(c.Context())
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "Flags not found or an error occurred",
		})
	}
	return c.Status(fiber.StatusOK).JSON(flags)
}

func (h *Handlers) GetFlagsFlagsId(c *fiber.Ctx, flagsId api.FlagsId) error {
	//TODO implement me
	detailsFromToken, err := h.getDetailsFromToken(c)
	isAdmin := false
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(api.Error{
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
		return c.Status(fiber.StatusOK).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}
	flag, err := h.Moderation.GetFlagByID(c.Context(), flagsId)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "An error occurred or Flag not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(flag)
}

func (h *Handlers) PostFlagsFlagsIdReview(c *fiber.Ctx, flagsId api.FlagsId) error {
	//TODO implement me
	detailsFromToken, err := h.getDetailsFromToken(c)
	isAdmin := false
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
	var status api.PostFlagsFlagsIdReviewJSONBody
	if err := c.BodyParser(&status); err != nil {
		return c.Status(fiber.StatusOK).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}
	err = h.Moderation.ReviewFlag(c.Context(), flagsId, string(*status.Status))
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "An error occurred or Flag not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON("Successful")
}

func (h *Handlers) PostFlagsSongId(c *fiber.Ctx, songId api.SongId) error {
	//TODO implement me
	detailsFromToken, err := h.getDetailsFromToken(c)
	isAdmin := false
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(api.Error{
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
		return c.Status(fiber.StatusOK).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}
	err = h.Song.FlagSong(c.Context(), songId)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "An error occured or song not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON("Successful")
}
