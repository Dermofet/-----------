package db

import (
	"context"
	"music-backend-test/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./source_mock.go -package=db
type UserSource interface {
	CreateUser(ctx context.Context, user *entity.UserCreate) (*entity.UserID, error)
	GetUserById(ctx context.Context, id *entity.UserID) (*entity.UserDB, error)
	GetUserByUsername(ctx context.Context, email string) (*entity.UserDB, error)
	UpdateUser(ctx context.Context, id *entity.UserID, user *entity.UserCreate) (*entity.UserDB, error)
	DeleteUser(ctx context.Context, id *entity.UserID) error
}

type MusicSource interface {
	GetAll(ctx context.Context) ([]*entity.Music, error)
	Create(ctx context.Context, musicCreate *entity.MusicCreate) error
	Update(ctx context.Context, id *entity.MusicID, musicUpdate *entity.MusicDB) error
	Delete(ctx context.Context, id *entity.MusicID) error
}
