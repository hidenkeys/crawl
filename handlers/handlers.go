package handlers

import (
	"crawl/repositories"
	"crawl/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Handlers struct {
	User       services.UserService
	Artist     services.ArtistService
	Album      services.AlbumService
	Song       services.SongService
	Genre      services.GenreService
	Playlist   services.PlaylistService
	Purchase   services.PurchaseService
	Stream     services.StreamService
	Tip        services.TipService
	Moderation services.ModerationService
	Auth       services.AuthService
	Payment    services.PaymentService
}

func (h *Handlers) PostStripeWebhook(c *fiber.Ctx) error {
	if err := h.Payment.HandleWebhook(c.Body(), c.Get("Stripe-Signature")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON("Successful")
	//panic("Implement me")
}

func NewHandlers(db *gorm.DB) *Handlers {
	repos := repositories.NewRepositories(db)
	return &Handlers{
		User:       services.NewUserService(repos.User, repos.Role, repos.Playlist, repos.Artist, repos.SongPurchase, repos.AlbumPurchase, repos.Stream),
		Artist:     services.NewArtistService(repos.Artist, repos.Song, repos.User, repos.Role),
		Album:      services.NewAlbumService(repos.Album, repos.AlbumContributor, repos.Song, repos.Artist),
		Song:       services.NewSongService(repos.Song, repos.Artist, repos.Genre, repos.Album, repos.Stream, repos.SongContributorRepository),
		Genre:      services.NewGenreService(repos.Genre),
		Playlist:   services.NewPlaylistService(repos.Playlist, repos.PlaylistSong, repos.Song),
		Purchase:   services.NewPurchaseService(repos.AlbumPurchase, repos.SongPurchase, repos.Album, repos.Song),
		Stream:     services.NewStreamService(repos.Stream, repos.Song),
		Tip:        services.NewTipService(repos.Tip, repos.User, repos.Artist),
		Moderation: services.NewModerationService(repos.Moderation),
		Auth:       services.NewAuthService(repos.User, repos.Role, repos.Artist),
		Payment:    services.NewPaymentService(repos.AlbumPurchase, repos.SongPurchase, repos.Album, repos.Song),
	}
}
