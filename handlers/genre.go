package handlers

import (
	"crawl/api"
	"crawl/models"
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

func (h *Handlers) PostGenres(c *fiber.Ctx) error {
	detailsFromToken, err := h.getDetailsFromToken(c)
	isAdmin := false
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
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
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var genreReq api.PostGenresJSONRequestBody
	if err := c.BodyParser(&genreReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	newGenre := &models.Genre{
		Name:        genreReq.Name,
		Description: genreReq.Description,
	}
	if genreReq.ImageUrl != nil {
		newGenre.ImageURL = *genreReq.ImageUrl
	}

	genre, err := h.Genre.Create(c.Context(), newGenre)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "An error occurred creating genre",
		})
	}
	return c.Status(fiber.StatusOK).JSON(genre)
}

func (h *Handlers) PutGenresGenreId(c *fiber.Ctx, genreId api.GenreId) error {
	detailsFromToken, err := h.getDetailsFromToken(c)
	isAdmin := false
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
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
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var genreReq api.PutGenresGenreIdJSONRequestBody
	if err := c.BodyParser(&genreReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	genre, err := h.Genre.GetGenreByID(c.Context(), genreId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "Genre not found for given ID",
		})
	}

	genre.Name = genreReq.Name
	genre.Description = genreReq.Description
	if genreReq.ImageUrl != nil {
		genre.ImageURL = *genreReq.ImageUrl
	}

	updatedGenre, err := h.Genre.Update(c.Context(), genre)
	if err != nil {
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to update genre",
		})
	}
	return c.Status(fiber.StatusOK).JSON(updatedGenre)
}
