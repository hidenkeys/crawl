package handlers

import (
	"crawl/api"
	"crawl/models"
	"crawl/services"
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) PostPurchasesAlbums(c *fiber.Ctx) error {
	userDetails, err := h.getDetailsFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var purchaseReq api.PostPurchasesAlbumsJSONBody
	if err := c.BodyParser(&purchaseReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	// Verify the requesting user is purchasing for themselves
	if userDetails.userID != purchaseReq.UserId {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "You can only purchase for yourself",
		})
	}

	// Verify the album exists
	album, err := h.Album.GetAlbumByID(c.Context(), purchaseReq.AlbumId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Album not found",
		})
	}

	// Process payment
	purchase, err := h.Purchase.PurchaseAlbum(c.Context(), purchaseReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(api.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed to process purchase",
		})
	}
	checkoutReq := services.CreateCheckoutSessionRequest{
		Name:     album.Title,
		ID:       purchase.ID.String(),
		Price:    int64(album.Price),
		ItemType: "album",
		ItemID:   album.ID.String(),
		UserID:   userDetails.userID.String(),
	}

	session, err := h.Payment.CreateCheckoutSession(checkoutReq)
	if err != nil {
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to create checkout session",
		})
	}

	purchase.StripeTransactionID = session.StripeSessionID
	purchase.Metadata = session.MetaData

	_, err = h.Purchase.UpdatePurchaseAlbum(c.Context(), *purchase)
	if err != nil {
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to update purchase",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "Album purchase request made successfully",
		Data: map[string]string{
			"url": session.PaymentUrl,
		},
	})
}

func (h *Handlers) PostPurchasesSongs(c *fiber.Ctx) error {
	userDetails, err := h.getDetailsFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(api.Error{
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	var purchaseReq api.PostPurchasesSongsJSONBody
	if err := c.BodyParser(&purchaseReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
		})
	}

	// Verify the requesting user is purchasing for themselves
	if userDetails.userID != purchaseReq.UserId {
		return c.Status(fiber.StatusForbidden).JSON(api.Error{
			Code:    fiber.StatusForbidden,
			Message: "Invalid user ID",
		})
	}

	// Verify the song exists
	song, err := h.Song.GetSongByID(c.Context(), purchaseReq.SongId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(api.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Song not found",
		})
	}

	// Process payment
	purchase, err := h.Purchase.PurchaseSong(c.Context(), purchaseReq)
	if err != nil {
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to process purchase",
		})
	}

	checkoutReq := services.CreateCheckoutSessionRequest{
		Name:     song.Title,
		ID:       purchase.ID.String(),
		Price:    int64(song.Price),
		ItemType: "song",
		ItemID:   song.ID.String(),
		UserID:   userDetails.userID.String(),
	}

	session, err := h.Payment.CreateCheckoutSession(checkoutReq)
	if err != nil {
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to create checkout session",
		})
	}

	purchase.StripeTransactionID = session.StripeSessionID
	purchase.Metadata = session.MetaData

	_, err = h.Purchase.UpdatePurchaseSong(c.Context(), *purchase)
	if err != nil {
		return c.Status(fiber.StatusExpectationFailed).JSON(api.Error{
			Code:    fiber.StatusExpectationFailed,
			Message: "Failed to update purchase",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Code:    fiber.StatusOK,
		Message: "Song purchase request made successfully",
		Data: map[string]string{
			"url": session.PaymentUrl,
		},
	})
}
