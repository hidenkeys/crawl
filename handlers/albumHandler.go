package handlers

import (
	"crawl/api"
	"github.com/gofiber/fiber/v2"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (s Server) GetAllAlbums(c *fiber.Ctx, params api.GetAllAlbumsParams) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetPurchasedAlbums(c *fiber.Ctx, userId openapi_types.UUID, params api.GetPurchasedAlbumsParams) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) CreateAlbum(c *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) DeleteAlbum(c *fiber.Ctx, albumId openapi_types.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetAlbum(c *fiber.Ctx, albumId openapi_types.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) UpdateAlbum(c *fiber.Ctx, albumId openapi_types.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetSongsInAlbum(c *fiber.Ctx, albumId openapi_types.UUID, params api.GetSongsInAlbumParams) error {
	//TODO implement me
	panic("implement me")
}

func (s Server) GetArtistAlbums(c *fiber.Ctx, artistId openapi_types.UUID, params api.GetArtistAlbumsParams) error {
	//TODO implement me
	panic("implement me")
}
