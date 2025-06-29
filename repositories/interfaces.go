package repositories

import (
	"crawl/models"
	"github.com/google/uuid"
	"time"
)

// IBaseRepository Base operations
type IBaseRepository[T any] interface {
	GetAll(offset int, limit int, where ...interface{}) ([]T, error)
	GetByID(id uuid.UUID) (*T, error)
	Create(model *T) (*T, error)
	Update(model *T) (*T, error)
	Delete(id uuid.UUID) error
	Exists(id uuid.UUID) (bool, error)
}

// IUserRepository User operations
type IUserRepository interface {
	IBaseRepository[models.User]
	FindByEmail(email string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	GetWithRoles(id uuid.UUID) (*models.User, error)
	Search(query string) ([]models.User, error)
	SetUserAsArtist(userID uuid.UUID) error
}

// IArtistRepository Artist operations
type IArtistRepository interface {
	IBaseRepository[models.Artist]
	GetWithSongs(id uuid.UUID) (*models.Artist, error)
	GetWithAlbums(id uuid.UUID) (*models.Artist, error)
	GetWithUserId(userID uuid.UUID) (*models.Artist, error)
	SearchByName(query string, limit int, offset int) ([]models.Artist, error)
}

// ISongRepository song operations
type ISongRepository interface {
	IBaseRepository[models.Song]
	GetWithArtist(id uuid.UUID) (*models.Song, error)
	GetByArtist(artistID uuid.UUID) ([]models.Song, error)
	GetTrending(limit int, since time.Time) ([]models.Song, error)
	AddPlayCount(id uuid.UUID, count int) error
	Search(query, artist, genre *string, sort, order *string, offset, limit int) ([]models.Song, error)
}

type IAlbumRepository interface {
	IBaseRepository[models.Album]
	GetWithSongs(id uuid.UUID) (*models.Album, error)
	GetByArtist(artistID uuid.UUID) ([]models.Album, error)
	SearchAlbums(query *string, artist *string, genre *string, sort *string, page *int, limit *int) ([]models.Album, error)
}

type IPlaylistRepository interface {
	IBaseRepository[models.Playlist]
	GetWithSongs(id uuid.UUID) (*models.Playlist, error)
	GetUserPlaylists(userID uuid.UUID) ([]models.Playlist, error)
	AddSongToPlaylist(playlistID, songID uuid.UUID) error
	GetUserPublicPlaylists(userID uuid.UUID) ([]models.Playlist, error)
	SearchPlaylists(
		query *string,
		owner *string,
		isPublic *bool,
		sort *string,
		page int,
		limit int,
	) ([]models.Playlist, int64, error)
}
type IStreamRepository interface {
	IBaseRepository[models.Stream]
	GetStreamCount(songID uuid.UUID, since time.Time) (int64, error)
	GetArtistStreams(artistID uuid.UUID, start, end time.Time) ([]models.Stream, error)
	GetStreamBySong(songID uuid.UUID) (*models.Stream, error)
}

type ITipRepository interface {
	IBaseRepository[models.ArtistTip]
	GetArtistTips(artistID uuid.UUID) ([]models.ArtistTip, error)
	GetUserTips(userID uuid.UUID) ([]models.ArtistTip, error)
	GetTotalTipsReceived(artistID uuid.UUID) (int64, error)
	GetTotalTipsSent(userID uuid.UUID) (int64, error)
}

type IModerationRepository interface {
	IBaseRepository[models.ContentFlag]
	GetFlaggedContent() ([]models.ContentFlag, error)
	UpdateFlagStatus(id uuid.UUID, status string) error
}

type IGenreRepository interface {
	IBaseRepository[models.Genre]
	GetPopular(limit int) ([]models.Genre, error)
	SearchGenres(query *string, sort *string) ([]models.Genre, error)
}

// ISongContributorRepository Song Contributor
type ISongContributorRepository interface {
	FindBySongID(songID uuid.UUID) ([]models.SongContributor, error)
	FindByArtistID(artistID uuid.UUID) ([]models.SongContributor, error)
	AddContributor(contributor *models.SongContributor) error
	RemoveContributor(songID, artistID uuid.UUID, contributionType string) error
	UpdateRoyalty(songID, artistID uuid.UUID, contributionType string, royalty int) error
}

// IAlbumContributorRepository Album Contributor
type IAlbumContributorRepository interface {
	FindByAlbumID(albumID uuid.UUID) ([]models.AlbumContributor, error)
	AddContributor(contributor *models.AlbumContributor) error
	RemoveContributor(albumID, artistID uuid.UUID, contributionType string) error
}

// IRoleRepository Role
type IRoleRepository interface {
	IBaseRepository[models.Role]
	AssignRoleToUser(userID, roleID uuid.UUID) error
	RemoveRoleFromUser(userID, roleID uuid.UUID) error
	GetUserRoles(userID uuid.UUID) ([]models.Role, error)
	FindByRolename(name string) (*models.Role, error)
}

// IUserFavoriteRepository User Favorite
type IUserFavoriteRepository interface {
	AddFavorite(userID, songID uuid.UUID) error
	RemoveFavorite(userID, songID uuid.UUID) error
	GetUserFavorites(userID uuid.UUID) ([]models.Song, error)
	IsFavorite(userID, songID uuid.UUID) (bool, error)
}

// IPlaylistSongRepository Playlist Song
type IPlaylistSongRepository interface {
	AddSongToPlaylist(playlistID, songID uuid.UUID, position int) error
	RemoveSongFromPlaylist(playlistID, songID uuid.UUID) error
	ReorderSongs(playlistID uuid.UUID, songOrder map[uuid.UUID]int) error
	GetPlaylistSongs(playlistID uuid.UUID) ([]models.Song, error)
	SongInPlaylistExists(playlistID uuid.UUID, songID uuid.UUID) (bool, error)
}

// ISongPurchaseRepository Song Purchase
type ISongPurchaseRepository interface {
	IBaseRepository[models.SongPurchase]
	FindByUserAndSong(userID, songID uuid.UUID) (*models.SongPurchase, error)
	GetUserPurchases(userID uuid.UUID) ([]models.SongPurchase, error)
	GetUserSongPurchases(userID uuid.UUID) ([]models.SongPurchase, error)
	HasPurchasedSong(userID, songID uuid.UUID) (bool, error)
}

// IAlbumPurchaseRepository Album Purchase
type IAlbumPurchaseRepository interface {
	IBaseRepository[models.AlbumPurchase]
	FindByUserAndAlbum(userID, albumID uuid.UUID) (*models.AlbumPurchase, error)
	GetUserAlbumPurchases(userID uuid.UUID) ([]models.AlbumPurchase, error)
	HasPurchasedAlbum(userID, albumID uuid.UUID) (bool, error)
}

// IMonthlyRoyaltyRepository Monthly Royalty
type IMonthlyRoyaltyRepository interface {
	IBaseRepository[models.MonthlyRoyalty]
	FindByArtistAndPeriod(artistID uuid.UUID, year int, month int) (*models.MonthlyRoyalty, error)
	MarkAsPaid(artistID uuid.UUID, year int, month int) error
	GetArtistRoyalties(artistID uuid.UUID) ([]models.MonthlyRoyalty, error)
	CalculatePendingRoyalties() (float64, error)
}
