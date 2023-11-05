package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"music-backend-test/internal/entity"
	"music-backend-test/internal/repository"

	"github.com/google/uuid"
)

type userInteractor struct {
	repo repository.UserRepository
}

func NewUserInteractor(repo repository.UserRepository) *userInteractor {
	return &userInteractor{
		repo: repo,
	}
}

func (u *userInteractor) Create(ctx context.Context, user *entity.UserCreate) (uuid.UUID, error) {
	userId, err := u.repo.Create(ctx, user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("can't create user by repository: %w", err)
	}

	return userId, nil
}

func (u *userInteractor) GetById(ctx context.Context, id uuid.UUID) (*entity.UserDB, error) {
	user, err := u.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("can't get user by id from repository: %w", err)
	}

	return user, nil
}

func (u *userInteractor) GetByUsername(ctx context.Context, email string) (*entity.UserDB, error) {
	user, err := u.repo.GetByUsername(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("can't get user by email from repository %w", err)
	}

	return user, nil
}

func (u *userInteractor) Update(ctx context.Context, id uuid.UUID, user *entity.UserCreate) (*entity.UserDB, error) {
	dbUser, err := u.repo.Update(ctx, id, user)
	if err != nil {
		return nil, fmt.Errorf("can't update user by repository: %w", err)
	}

	return dbUser, nil
}

func (u *userInteractor) Delete(ctx context.Context, id uuid.UUID) error {
	err := u.repo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return fmt.Errorf("can't delete user by repository: %w", err)
	}

	return nil
}

func (u *userInteractor) LikeTrack(ctx context.Context, userId uuid.UUID, trackId uuid.UUID) error {
	err := u.repo.LikeTrack(ctx, userId, trackId)
	if err != nil {
		return fmt.Errorf("/repository/user.LikeTrack: %w", err)
	}

	return nil
}

func (u *userInteractor) DislikeTrack(ctx context.Context, userId uuid.UUID, trackId uuid.UUID) error {
	err := u.repo.DislikeTrack(ctx, userId, trackId)
	if err != nil {
		return fmt.Errorf("/repository/user.DislikeTrack: %w", err)
	}

	return nil
}

func (u *userInteractor) ShowLikedTracks(ctx context.Context, id uuid.UUID) ([]*entity.MusicDB, error) {
	data, err := u.repo.ShowLikedTracks(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("/repository/user.ShowLikedTracks: %w", err)
	}

	return data, nil
}
