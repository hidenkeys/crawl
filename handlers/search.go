package handlers

import (
	"crawl/api"
	"crawl/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type MediaItem struct {
	Type  string        `json:"type"` // "song" or "album"
	Song  *models.Song  `json:"song,omitempty"`
	Album *models.Album `json:"album,omitempty"`
}

type GlobalSearch struct {
	Media     []MediaItem       `json:"media"`
	Artists   []models.Artist   `json:"artists"`
	Genre     []models.Genre    `json:"genre"`
	Playlists []models.Playlist `json:"playlists"`
}

func (h *Handlers) GetSearch(c *fiber.Ctx, params api.GetSearchParams) error {
	songs, _ := h.Song.SearchSongs(c.Context(), &params.Query, nil, nil, nil, nil, params.Page, params.Limit)
	playlists, _, _ := h.Playlist.SearchPlaylists(c.Context(), &params.Query, nil, nil, nil, *params.Page, *params.Limit)
	genres, _ := h.Genre.SearchGenres(c.Context(), &params.Query, nil)
	artists, _ := h.Artist.SearchArtistsByName(c.Context(), params.Query, *params.Page, *params.Limit)
	albums, _ := h.Album.SearchAlbums(c.Context(), &params.Query, nil, nil, nil, params.Page, params.Limit)

	media := make([]MediaItem, 0, len(songs)+len(albums))
	for _, song := range songs {
		media = append(media, MediaItem{Type: "song", Song: &song})
	}
	for _, album := range albums {
		media = append(media, MediaItem{Type: "album", Album: &album})
	}

	result := &GlobalSearch{
		Media:     media,
		Artists:   artists,
		Genre:     genres,
		Playlists: playlists,
	}

	return c.JSON(result)
}

func (h *Handlers) GetSearchAlbums(c *fiber.Ctx, params api.GetSearchAlbumsParams) error {
	albums, err := h.Album.SearchAlbums(c.Context(), params.Query, params.Artist, params.Genre, (*string)(params.Sort), params.Page, params.Limit)
	if err != nil {
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to search albums",
		})
	}

	return c.JSON(albums)
}

func (h *Handlers) GetSearchArtists(c *fiber.Ctx, params api.GetSearchArtistsParams) error {
	artists, err := h.Artist.SearchArtistsByName(c.Context(), *params.Query, *params.Page, *params.Limit)
	if err != nil {
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to search artists",
		})
	}

	return c.JSON(artists)
}

func (h *Handlers) GetSearchGenres(c *fiber.Ctx, params api.GetSearchGenresParams) error {
	genres, err := h.Genre.SearchGenres(c.Context(), params.Query, (*string)(params.Sort))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to search genres",
		})
	}

	return c.JSON(genres)
}

func (h *Handlers) GetSearchPlaylists(c *fiber.Ctx, params api.GetSearchPlaylistsParams) error {
	playlists, _, err := h.Playlist.SearchPlaylists(c.Context(), params.Query, params.Owner, params.IsPublic, (*string)(params.Sort), *params.Page, *params.Limit)
	if err != nil {
		log.Infof("Error occured: %s", err.Error())
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to search playlists",
		})
	}

	return c.JSON(playlists)
}

func (h *Handlers) GetSearchSongs(c *fiber.Ctx, params api.GetSearchSongsParams) error {
	songs, err := h.Song.SearchSongs(c.Context(), params.Query, params.Artist, params.Genre, (*string)(params.Sort), (*string)(params.Order), params.Page, params.Limit)
	if err != nil {
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to search songs",
		})
	}

	return c.JSON(songs)
}
