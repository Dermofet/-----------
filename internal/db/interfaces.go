package db

import (
	"context"
	"music-backend-test/internal/entity"

	"github.com/google/uuid"
)

//go:generate mockgen -source=interfaces.go -destination=./source_mock.go -package=db
type UserSource interface {
	CreateUser(ctx context.Context, user *entity.UserCreate) (uuid.UUID, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*entity.UserDB, error)
	GetUserByUsername(ctx context.Context, email string) (*entity.UserDB, error)
	UpdateUser(ctx context.Context, userDB *entity.UserDB, user *entity.UserCreate) (*entity.UserDB, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	LikeTrack(ctx context.Context, userId uuid.UUID, trackId uuid.UUID) error
	DislikeTrack(ctx context.Context, userId uuid.UUID, trackId uuid.UUID) error
	ShowLikedTracks(ctx context.Context, id uuid.UUID) ([]*entity.MusicDB, error)
}

type MusicSource interface {
	GetAll(ctx context.Context) ([]*entity.MusicDB, error)
	Get(ctx context.Context, musicId uuid.UUID) (*entity.MusicDB, error)
	GetAndSortByPopular(ctx context.Context) ([]*entity.MusicDB, error)
	GetAllSortByTime(ctx context.Context) ([]*entity.MusicDB, error)
	Create(ctx context.Context, musicDb *entity.MusicDB) error
	Update(ctx context.Context, musicDb *entity.MusicDB) error
	Delete(ctx context.Context, id uuid.UUID) error
}
