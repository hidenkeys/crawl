package handlers

import (
	"crawl/api"
	"github.com/gofiber/fiber/v2"
	"github.com/oapi-codegen/runtime/types"
)

func (h *Handlers) GetPlaylistsPlaylistId(c *fiber.Ctx, playlistId types.UUID) error {
	playlist, err := h.Playlist.GetPlaylistByID(c.Context(), playlistId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "Playlist not found",
		})
	}

	// Check if playlist is private
	if !playlist.IsPublic {
		userID, err := h.getUserIDFromToken(c)
		if err != nil || userID != playlist.UserID {
			return c.Status(fiber.StatusForbidden).JSON(api.Error{
				Code:    fiber.StatusForbidden,
				Message: "You don't have access to this playlist",
			})
		}
	}

	return c.JSON(playlist)
}

func (h *Handlers) PostPlaylistsPlaylistIdSongs(c *fiber.Ctx, playlistId types.UUID) error {
	userID, err := h.getUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var songReq api.PostPlaylistsPlaylistIdSongsJSONBody
	if err := c.BodyParser(&songReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	// Verify the requesting user owns the playlist
	playlist, err := h.Playlist.GetPlaylistByID(c.Context(), playlistId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "Playlist not found",
		})
	}

	if playlist.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can only add songs to your own playlists",
		})
	}

	// Verify the song exists
	_, err = h.Song.GetSongByID(c.Context(), songReq.SongId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Song not found",
		})
	}

	if err := h.Playlist.AddSongToPlaylist(c.Context(), playlistId, songReq.SongId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to add song to playlist",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handlers) DeletePlaylistsPlaylistIdSongsSongId(c *fiber.Ctx, playlistId types.UUID, songId types.UUID) error {
	userID, err := h.getUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	// Verify the requesting user owns the playlist
	playlist, err := h.Playlist.GetPlaylistByID(c.Context(), playlistId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "Playlist not found",
		})
	}

	if playlist.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can only remove songs from your own playlists",
		})
	}

	if err := h.Playlist.RemoveSongFromPlaylist(c.Context(), playlistId, songId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to remove song from playlist",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handlers) GetPlaylistsPlaylistIdSongs(c *fiber.Ctx, playlistId types.UUID) error {
	playlist, err := h.Playlist.GetPlaylistByID(c.Context(), playlistId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "Playlist not found",
		})
	}

	// Check if playlist is private
	if !playlist.IsPublic {
		userID, err := h.getUserIDFromToken(c)
		if err != nil || userID != playlist.UserID {
			return c.Status(fiber.StatusForbidden).JSON(api.Error{
				Code:    fiber.StatusForbidden,
				Message: "You don't have access to this playlist",
			})
		}
	}

	songs, err := h.Playlist.GetPlaylistSongs(c.Context(), playlistId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to fetch playlist songs",
		})
	}

	return c.JSON(songs)
}

func (h *Handlers) PutPlaylistsPlaylistId(c *fiber.Ctx, playlistId types.UUID) error {
	userID, err := h.getUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var playlistReq api.Playlist
	if err := c.BodyParser(&playlistReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	// Verify the requesting user owns the playlist
	playlist, err := h.Playlist.GetPlaylistByID(c.Context(), playlistId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "Playlist not found",
		})
	}

	if playlist.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can only update your own playlists",
		})
	}

	playlist.Title = playlistReq.Title

	// Handle nullable/optional fields
	if playlistReq.Description != nil {
		playlist.Description = *playlistReq.Description
	}

	if playlistReq.CoverImageUrl != nil {
		playlist.CoverImageURL = *playlistReq.CoverImageUrl
	}

	if playlistReq.IsPublic != nil {
		playlist.IsPublic = *playlistReq.IsPublic
	} else {
		playlist.IsPublic = false // Default from gorm tag
	}

	updatedPlaylist, err := h.Playlist.UpdatePlaylist(c.Context(), playlistId, playlist)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to update playlist",
		})
	}

	return c.JSON(updatedPlaylist)
}

func (h *Handlers) DeletePlaylistsPlaylistId(c *fiber.Ctx, playlistId types.UUID) error {
	userID, err := h.getUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	// Verify the requesting user owns the playlist
	playlist, err := h.Playlist.GetPlaylistByID(c.Context(), playlistId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(api.Error{
			Code:    fiber.StatusNotFound,
			Message: "Playlist not found",
		})
	}

	if playlist.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can only delete your own playlists",
		})
	}

	if err := h.Playlist.DeletePlaylist(c.Context(), playlistId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to delete playlist",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
