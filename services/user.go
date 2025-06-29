package services

import (
	"context"
	"crawl/api"
	"crawl/models"
	"crawl/repositories"
	"errors"
	"github.com/google/uuid"
)

type UserService interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GetAllUsers(ctx context.Context, page int, limit int) ([]models.User, error)
	Search(ctx context.Context, query string) ([]models.User, error)
	Update(ctx context.Context, userID uuid.UUID, userReq api.User) (*models.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetArtistByUserId(ctx context.Context, userID uuid.UUID) (*models.Artist, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserPurchasedAlbums(ctx context.Context, userID uuid.UUID) ([]models.Album, error)
	GetUserPurchaseHistory(ctx context.Context, userID uuid.UUID, params api.GetUsersUserIdLibraryPurchasesParams) ([]models.SongPurchase, error)
	GetUserPurchasedSongs(ctx context.Context, userID uuid.UUID) ([]models.Song, error)
	GetUserPublicPlaylists(ctx context.Context, userID uuid.UUID) ([]models.Playlist, error)
	GetUserPlaylists(ctx context.Context, userID uuid.UUID) ([]models.Playlist, error)
	CreatePlaylist(ctx context.Context, userID uuid.UUID, playlist *models.Playlist) (*models.Playlist, error)
}

type userService struct {
	userRepo          repositories.IUserRepository
	roleRepo          repositories.IRoleRepository
	playlistRepo      repositories.IPlaylistRepository
	artistRepo        repositories.IArtistRepository
	songPurchaseRepo  repositories.ISongPurchaseRepository
	albumPurchaseRepo repositories.IAlbumPurchaseRepository
}

func NewUserService(
	userRepo repositories.IUserRepository,
	roleRepo repositories.IRoleRepository,
	playlistRepo repositories.IPlaylistRepository,
	artistRepo repositories.IArtistRepository,
	songPurchaseRepo repositories.ISongPurchaseRepository,
	albumPurchaseRepo repositories.IAlbumPurchaseRepository,

) UserService {
	return &userService{
		userRepo:          userRepo,
		roleRepo:          roleRepo,
		playlistRepo:      playlistRepo,
		artistRepo:        artistRepo,
		songPurchaseRepo:  songPurchaseRepo,
		albumPurchaseRepo: albumPurchaseRepo,
	}
}

func (s *userService) Create(ctx context.Context, user *models.User) (*models.User, error) {
	// Check if email already exists
	_, err := s.userRepo.FindByEmail(user.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	// Check if username already exists
	_, err = s.userRepo.FindByUsername(user.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	}

	role, err := s.roleRepo.FindByRolename("Listener")
	if err != nil {
		return nil, errors.New("role not found")
	}
	newUser, err := s.userRepo.Create(user)
	if err != nil {
		return nil, errors.New("error occurred creating user")
	}

	err = s.roleRepo.AssignRoleToUser(role.ID, newUser.ID)
	if err != nil {
		return nil, errors.New("error occurred adding Role to user")
	}

	return newUser, nil

}

func (s *userService) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(userID)
}

func (s *userService) GetAllUsers(ctx context.Context, page int, limit int) ([]models.User, error) {
	offset := (page - 1) * limit
	return s.userRepo.GetAll(offset, limit)
}

func (s *userService) Search(ctx context.Context, query string) ([]models.User, error) {
	return s.userRepo.Search(query)
}

func (s *userService) Update(ctx context.Context, userID uuid.UUID, userReq api.User) (*models.User, error) {
	// Implementation depends on what fields you want to update
	// This is just a placeholder
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	user.Email = string(userReq.Email)
	user.Username = userReq.Username
	user.FirstName = userReq.FirstName
	user.LastName = userReq.LastName
	if userReq.Bio != nil {
		user.Bio = *userReq.Bio
	}
	if userReq.PhoneNumber != nil {
		user.PhoneNumber = *userReq.PhoneNumber
	}
	if userReq.ProfileImageUrl != nil {
		user.ProfileImage = *userReq.ProfileImageUrl
	}
	return s.userRepo.Update(user)
}

func (s *userService) Delete(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.Delete(userID)
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.userRepo.FindByUsername(username)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.FindByEmail(email)
}

func (s *userService) GetUserPurchasedAlbums(ctx context.Context, userID uuid.UUID) ([]models.Album, error) {
	purchases, err := s.albumPurchaseRepo.GetUserAlbumPurchases(userID)
	if err != nil {
		return nil, err
	}

	albums := make([]models.Album, len(purchases))
	for i, p := range purchases {
		albums[i] = p.Album
	}

	return albums, nil
}

func (s *userService) GetUserPurchaseHistory(ctx context.Context, userID uuid.UUID, params api.GetUsersUserIdLibraryPurchasesParams) ([]models.SongPurchase, error) {
	return s.songPurchaseRepo.GetUserSongPurchases(userID)
}

func (s *userService) GetUserPurchasedSongs(ctx context.Context, userID uuid.UUID) ([]models.Song, error) {
	purchases, err := s.songPurchaseRepo.GetUserSongPurchases(userID)
	if err != nil {
		return nil, err
	}

	songs := make([]models.Song, len(purchases))
	for i, p := range purchases {
		songs[i] = p.Song
	}

	return songs, nil
}

func (s *userService) GetUserPlaylists(ctx context.Context, userID uuid.UUID) ([]models.Playlist, error) {
	return s.playlistRepo.GetUserPlaylists(userID)
}

func (s *userService) GetUserPublicPlaylists(ctx context.Context, userID uuid.UUID) ([]models.Playlist, error) {
	return s.playlistRepo.GetUserPublicPlaylists(userID)
}

func (s *userService) CreatePlaylist(ctx context.Context, userID uuid.UUID, playlist *models.Playlist) (*models.Playlist, error) {
	playlist.UserID = userID
	return s.playlistRepo.Create(playlist)
}

func (s *userService) GetArtistByUserId(ctx context.Context, userID uuid.UUID) (*models.Artist, error) {
	return s.artistRepo.GetWithUserId(userID)
}
