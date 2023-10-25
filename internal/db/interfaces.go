package db

import (
	"context"
	"go-test-grpc-http/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./source_mock.go -package=db

type UserSource interface {
	CreateUser(ctx context.Context, user *entity.UserCreate) (*entity.UserID, error)
	GetUserById(ctx context.Context, id *entity.UserID) (*entity.UserDB, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.UserDB, error)
	GetUserIdByEmail(ctx context.Context, email string) (*entity.UserID, error)
	UpdateUser(ctx context.Context, id *entity.UserID, user *entity.UserCreate) (*entity.UserDB, error)
	DeleteUser(ctx context.Context, id *entity.UserID) error
}
