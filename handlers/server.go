package handlers

import (
	"crawl/services"
	"gorm.io/gorm"
)

type Server struct {
	db                 *gorm.DB
	usrService         *services.UserService
	albumService       *services.AlbumService
	purchaseService    *services.PurchaseService
	songService        *services.SongService
	songMetricsService *services.SongMetricsService
	tipService         *services.TipService
}

func NewServer(db *gorm.DB, usrService *services.UserService,
	albumService *services.AlbumService,
	purchaseService *services.PurchaseService,
	songService *services.SongService,
	songMetricsService *services.SongMetricsService,
	tipService *services.TipService,
) *Server {
	return &Server{
		db:                 db,
		usrService:         usrService,
		albumService:       albumService,
		purchaseService:    purchaseService,
		songService:        songService,
		songMetricsService: songMetricsService,
		tipService:         tipService,
	}
}
