package repository

import (
	"context"
	"go-test-grpc-http/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./repositories_mock.go -package=repository

type UserRepository interface {
	Create(ctx context.Context, user *entity.UserCreate) (*entity.UserID, error)
	GetById(ctx context.Context, id *entity.UserID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetIdByEmail(ctx context.Context, email string) (*entity.UserID, error)
	Update(ctx context.Context, id *entity.UserID, user *entity.UserCreate) (*entity.User, error)
	Delete(ctx context.Context, id *entity.UserID) error
}
