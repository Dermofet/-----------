package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"go-test-grpc-http/internal/entity"
	"go-test-grpc-http/internal/repository"
)

type userInteractor struct {
	repo repository.UserRepository
}

func NewUserInteractor(repo repository.UserRepository) *userInteractor {
	return &userInteractor{
		repo: repo,
	}
}

func (u *userInteractor) Create(ctx context.Context, user *entity.UserCreate) (*entity.UserID, error) {
	userId, err := u.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("can't create user by repository: %w", err)
	}

	return userId, nil
}

func (u *userInteractor) GetById(ctx context.Context, id *entity.UserID) (*entity.User, error) {
	user, err := u.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("can't get user by id from repository: %w", err)
	}

	return user, nil
}

func (u *userInteractor) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("can't get user by email from repository %w", err)
	}

	return user, nil
}

func (u *userInteractor) GetIdByEmail(ctx context.Context, email string) (*entity.UserID, error) {
	id, err := u.repo.GetIdByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("can't get user id by email from repository %w", err)
	}

	return id, nil
}

func (u *userInteractor) Update(ctx context.Context, id *entity.UserID, user *entity.UserCreate) (*entity.User, error) {
	dbUser, err := u.repo.Update(ctx, id, user)
	if err != nil {
		return nil, fmt.Errorf("can't update user by repository: %w", err)
	}

	return dbUser, nil
}

func (u *userInteractor) Delete(ctx context.Context, id *entity.UserID) error {
	err := u.repo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return fmt.Errorf("can't delete user by repository: %w", err)
	}

	return nil
}
