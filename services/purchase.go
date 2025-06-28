package services

import (
	"context"
	"crawl/api"
	"crawl/models"
	"crawl/repositories"
	"errors"
	"github.com/google/uuid"
)

type PurchaseService interface {
	PurchaseAlbum(ctx context.Context, purchase api.PostPurchasesAlbumsJSONBody) (*models.AlbumPurchase, error)
	UpdatePurchaseAlbum(ctx context.Context, albumPurchase models.AlbumPurchase) (*models.AlbumPurchase, error)
	GetAllPurchaseAlbumByUser(ctx context.Context, userID uuid.UUID) ([]models.AlbumPurchase, error)
	GetPurchaseAlbumByIdByUser(ctx context.Context, purchaseAlbumID uuid.UUID, userID uuid.UUID) (*models.AlbumPurchase, error)
	GetAllPurchaseAlbum(ctx context.Context, userID uuid.UUID) ([]models.AlbumPurchase, error)
	GetPurchaseAlbumById(ctx context.Context, albumPurchaseID uuid.UUID) (*models.AlbumPurchase, error)
	GetPurchasedAlbum(ctx context.Context, userID uuid.UUID) ([]models.Album, error)

	PurchaseSong(ctx context.Context, purchase api.PostPurchasesSongsJSONBody) (*models.SongPurchase, error)
	UpdatePurchaseSong(ctx context.Context, songPurchase models.SongPurchase) (*models.SongPurchase, error)
	GetAllPurchaseSongs(ctx context.Context, userID uuid.UUID) ([]models.SongPurchase, error)
	GetPurchaseSongById(ctx context.Context, songPurchaseID uuid.UUID) (*models.SongPurchase, error)
	GetAllPurchaseSongByUser(ctx context.Context, userID uuid.UUID) ([]models.SongPurchase, error)
	GetPurchaseSongByIdByUser(ctx context.Context, purchaseSongID uuid.UUID, userID uuid.UUID) (*models.SongPurchase, error)
	GetPurchasedSong(ctx context.Context, userID uuid.UUID) ([]models.Song, error)
}

type purchaseService struct {
	albumPurchaseRepo repositories.IAlbumPurchaseRepository
	songPurchaseRepo  repositories.ISongPurchaseRepository
	albumRepo         repositories.IAlbumRepository
	songRepo          repositories.ISongRepository
}

func NewPurchaseService(
	albumPurchaseRepo repositories.IAlbumPurchaseRepository,
	songPurchaseRepo repositories.ISongPurchaseRepository,
	albumRepo repositories.IAlbumRepository,
	songRepo repositories.ISongRepository,
) PurchaseService {
	return &purchaseService{
		albumPurchaseRepo: albumPurchaseRepo,
		songPurchaseRepo:  songPurchaseRepo,
		albumRepo:         albumRepo,
		songRepo:          songRepo,
	}
}

// Album Purchase Methods

func (s *purchaseService) PurchaseAlbum(ctx context.Context, purchase api.PostPurchasesAlbumsJSONBody) (*models.AlbumPurchase, error) {
	// Check if album exists
	album, err := s.albumRepo.GetByID(purchase.AlbumId)
	if err != nil {
		return nil, errors.New("album not found")
	}

	// Check if user already purchased this album
	hasPurchased, err := s.albumPurchaseRepo.HasPurchasedAlbum(purchase.UserId, purchase.AlbumId)
	if err != nil {
		return nil, err
	}
	if hasPurchased {
		return nil, errors.New("user already purchased this album")
	}

	// Create new purchase
	newPurchase := &models.AlbumPurchase{
		UserID:              purchase.UserId,
		AlbumID:             purchase.AlbumId,
		PurchasePrice:       float64(album.Price),
		Currency:            "NGN",
		PaymentStatus:       "completed",
		StripeTransactionID: purchase.PaymentMethodId,
	}

	return s.albumPurchaseRepo.Create(newPurchase)
}

func (s *purchaseService) UpdatePurchaseAlbum(ctx context.Context, albumPurchase models.AlbumPurchase) (*models.AlbumPurchase, error) {
	// Verify purchase exists
	existingPurchase, err := s.albumPurchaseRepo.GetByID(albumPurchase.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return nil, errors.New("purchase not found")
		}
		return nil, err
	}

	// Update fields
	existingPurchase.PaymentStatus = albumPurchase.PaymentStatus
	existingPurchase.StripeTransactionID = albumPurchase.StripeTransactionID

	return s.albumPurchaseRepo.Update(existingPurchase)
}

func (s *purchaseService) GetAllPurchaseAlbumByUser(ctx context.Context, userID uuid.UUID) ([]models.AlbumPurchase, error) {
	return s.albumPurchaseRepo.GetUserAlbumPurchases(userID)
}

func (s *purchaseService) GetPurchaseAlbumByIdByUser(ctx context.Context, purchaseAlbumID uuid.UUID, userID uuid.UUID) (*models.AlbumPurchase, error) {
	purchase, err := s.albumPurchaseRepo.GetByID(purchaseAlbumID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return nil, errors.New("purchase not found")
		}
		return nil, err
	}

	if purchase.UserID != userID {
		return nil, errors.New("unauthorized access to purchase")
	}

	return purchase, nil
}

func (s *purchaseService) GetAllPurchaseAlbum(ctx context.Context, userID uuid.UUID) ([]models.AlbumPurchase, error) {
	return s.albumPurchaseRepo.GetUserAlbumPurchases(userID)
}

func (s *purchaseService) GetPurchaseAlbumById(ctx context.Context, albumPurchaseID uuid.UUID) (*models.AlbumPurchase, error) {
	purchase, err := s.albumPurchaseRepo.GetByID(albumPurchaseID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return nil, errors.New("purchase not found")
		}
		return nil, err
	}
	return purchase, nil
}

func (s *purchaseService) GetPurchasedAlbum(ctx context.Context, userID uuid.UUID) ([]models.Album, error) {
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

// Song Purchase Methods

func (s *purchaseService) PurchaseSong(ctx context.Context, purchase api.PostPurchasesSongsJSONBody) (*models.SongPurchase, error) {
	// Check if song exists
	var song *models.Song
	song, err := s.songRepo.GetByID(purchase.SongId)
	if err != nil {
		return nil, errors.New("song not found")
	}

	// Check if user already purchased this song
	hasPurchased, err := s.songPurchaseRepo.HasPurchasedSong(purchase.UserId, purchase.SongId)
	if err != nil {
		return nil, err
	}
	if hasPurchased {
		return nil, errors.New("user already purchased this song")
	}

	// Create new purchase
	newPurchase := &models.SongPurchase{
		UserID:              purchase.UserId,
		SongID:              purchase.SongId,
		PurchasePrice:       float64(song.Price),
		Currency:            "NGN",
		PaymentStatus:       "completed",
		StripeTransactionID: purchase.PaymentMethodId,
	}

	return s.songPurchaseRepo.Create(newPurchase)
}

func (s *purchaseService) UpdatePurchaseSong(ctx context.Context, songPurchase models.SongPurchase) (*models.SongPurchase, error) {
	// Verify purchase exists
	existingPurchase, err := s.songPurchaseRepo.GetByID(songPurchase.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return nil, errors.New("purchase not found")
		}
		return nil, err
	}

	// Update fields
	existingPurchase.PaymentStatus = songPurchase.PaymentStatus
	existingPurchase.StripeTransactionID = songPurchase.StripeTransactionID

	return s.songPurchaseRepo.Update(existingPurchase)
}

func (s *purchaseService) GetAllPurchaseSongs(ctx context.Context, userID uuid.UUID) ([]models.SongPurchase, error) {
	return s.songPurchaseRepo.GetUserSongPurchases(userID)
}

func (s *purchaseService) GetPurchaseSongById(ctx context.Context, songPurchaseID uuid.UUID) (*models.SongPurchase, error) {
	purchase, err := s.songPurchaseRepo.GetByID(songPurchaseID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return nil, errors.New("purchase not found")
		}
		return nil, err
	}
	return purchase, nil
}

func (s *purchaseService) GetAllPurchaseSongByUser(ctx context.Context, userID uuid.UUID) ([]models.SongPurchase, error) {
	return s.songPurchaseRepo.GetUserSongPurchases(userID)
}

func (s *purchaseService) GetPurchaseSongByIdByUser(ctx context.Context, purchaseSongID uuid.UUID, userID uuid.UUID) (*models.SongPurchase, error) {
	purchase, err := s.songPurchaseRepo.GetByID(purchaseSongID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return nil, errors.New("purchase not found")
		}
		return nil, err
	}

	if purchase.UserID != userID {
		return nil, errors.New("unauthorized access to purchase")
	}

	return purchase, nil
}

func (s *purchaseService) GetPurchasedSong(ctx context.Context, userID uuid.UUID) ([]models.Song, error) {
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
