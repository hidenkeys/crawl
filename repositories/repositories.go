package repositories

import (
	"gorm.io/gorm"
)

type Repositories struct {
	User                      IUserRepository
	Artist                    IArtistRepository
	Album                     IAlbumRepository
	Song                      ISongRepository
	Genre                     IGenreRepository
	Playlist                  IPlaylistRepository
	SongPurchase              ISongPurchaseRepository
	Stream                    IStreamRepository
	Tip                       ITipRepository
	Moderation                IModerationRepository
	AlbumPurchase             IAlbumPurchaseRepository
	Role                      IRoleRepository
	AlbumContributor          IAlbumContributorRepository
	SongContributorRepository ISongContributorRepository
	MonthlyRoyalty            IMonthlyRoyaltyRepository
	UserFavorite              IUserFavoriteRepository
	PlaylistSong              IPlaylistSongRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:                      NewUserRepository(db),
		Artist:                    NewArtistRepository(db),
		Album:                     NewAlbumRepository(db),
		Song:                      NewSongRepository(db),
		Genre:                     NewGenreRepository(db),
		Playlist:                  NewPlaylistRepository(db),
		SongPurchase:              NewSongPurchaseRepository(db),
		Stream:                    NewStreamRepository(db),
		Tip:                       NewTipRepository(db),
		Moderation:                NewModerationRepository(db),
		AlbumPurchase:             NewAlbumPurchaseRepository(db),
		Role:                      NewRoleRepository(db),
		AlbumContributor:          NewAlbumContributorRepository(db),
		SongContributorRepository: NewSongContributorRepository(db),
		MonthlyRoyalty:            NewMonthlyRoyaltyRepository(db),
		UserFavorite:              NewUserFavoriteRepository(db),
		PlaylistSong:              NewPlaylistSongRepository(db),
	}
}
