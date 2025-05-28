package handlers

import (
	"crawl/api"
	"github.com/gofiber/fiber/v2"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (s Server) GetPurchasedSongs(c *fiber.Ctx, userId openapi_types.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetAllSongs(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) CreateSong(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) SearchSongs(c *fiber.Ctx, params api.SearchSongsParams) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) DeleteSong(c *fiber.Ctx, songId openapi_types.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetSong(c *fiber.Ctx, songId openapi_types.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetArtistSongs(c *fiber.Ctx, artistId openapi_types.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) UpdateSong(c *fiber.Ctx, songId openapi_types.UUID) error {
	//TODO implement me
	panic("implement me")
}
