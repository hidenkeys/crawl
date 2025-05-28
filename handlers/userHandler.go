package handlers

import (
	"crawl/api"
	"github.com/gofiber/fiber/v2"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"net/http"
)

func (s Server) GetArtistDashboard(c *fiber.Ctx, artistId openapi_types.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetArtistRevenueDashboard(c *fiber.Ctx, artistId openapi_types.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) UserLogin(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) CreateUser(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")

	var reqBody api.CreateUserJSONRequestBody

	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(api.Error{
			ErrorCode: "400",
			Message:   "Invalid request body",
		})
	}
}

func (s Server) DeleteUser(c *fiber.Ctx, userId openapi_types.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetUser(c *fiber.Ctx, userId openapi_types.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) UpdateUser(c *fiber.Ctx, userId openapi_types.UUID) error {
	//TODO implement me
	panic("implement me")
}
