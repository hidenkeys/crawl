package handlers

import (
	"crawl/api"
	"github.com/gofiber/fiber/v2"
	"github.com/oapi-codegen/runtime/types"
)

func (h *Handlers) GetGenres(c *fiber.Ctx) error {
	genres, err := h.Genre.GetAllGenres(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to fetch genres",
		})
	}

	return c.JSON(genres)
}

func (h *Handlers) GetGenresGenreId(c *fiber.Ctx, genreId types.UUID) error {
	genre, err := h.Genre.GetGenreByID(c.Context(), genreId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "Genre not found",
		})
	}

	return c.JSON(genre)
}
