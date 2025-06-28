package handlers

import (
	"crawl/repositories"
	"crawl/services"
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
}

func NewHandlers(db *gorm.DB) *Handlers {
	repos := repositories.NewRepositories(db)
	return &Handlers{
		User:       services.NewUserService(repos.User, repos.Playlist, repos.Artist, repos.SongPurchase, repos.AlbumPurchase),
		Artist:     services.NewArtistService(repos.Artist, repos.Song, repos.User),
		Album:      services.NewAlbumService(repos.Album, repos.AlbumContributor, repos.Song),
		Song:       services.NewSongService(repos.Song, repos.Artist, repos.Genre, repos.Album, repos.Stream, repos.SongContributorRepository),
		Genre:      services.NewGenreService(repos.Genre),
		Playlist:   services.NewPlaylistService(repos.Playlist, repos.PlaylistSong, repos.Song),
		Purchase:   services.NewPurchaseService(repos.AlbumPurchase, repos.SongPurchase, repos.Album, repos.Song),
		Stream:     services.NewStreamService(repos.Stream, repos.Song),
		Tip:        services.NewTipService(repos.Tip, repos.User, repos.Artist),
		Moderation: services.NewModerationService(repos.Moderation),
		Auth:       services.NewAuthService(repos.User),
	}
}
