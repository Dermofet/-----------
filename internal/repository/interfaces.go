package repository

import (
	"context"
	"music-backend-test/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./repositories_mock.go -package=repository
type UserRepository interface {
	Create(ctx context.Context, user *entity.UserCreate) (*entity.UserID, error)
	GetById(ctx context.Context, id *entity.UserID) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	Update(ctx context.Context, id *entity.UserID, user *entity.UserCreate) (*entity.User, error)
	Delete(ctx context.Context, id *entity.UserID) error
	LikeTrack(ctx context.Context, userId *entity.UserID, trackId *entity.MusicID) error
	DislikeTrack(ctx context.Context, userId *entity.UserID, trackId *entity.MusicID) error
	ShowLikedTracks(ctx context.Context, id *entity.UserID) ([]*entity.Music, error)
}

type MusicRepository interface {
	GetAll(ctx context.Context) ([]*entity.Music, error)
	GetAndSortByPopular(ctx context.Context) ([]*entity.Music, error)
	GetAllSortByTime(ctx context.Context) ([]entity.MusicShow, error)
	Create(ctx context.Context, musicCreate *entity.MusicCreate) error
	Update(ctx context.Context, musicUpdate *entity.MusicDB) error
	Delete(ctx context.Context, id *entity.MusicID) error
}
