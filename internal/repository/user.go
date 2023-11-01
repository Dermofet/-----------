package repository

import (
	"context"
	"database/sql"
	"fmt"
	"music-backend-test/internal/db"
	"music-backend-test/internal/entity"
)

type userRepository struct {
	source db.UserSource
}

func NewUserRepository(source db.UserSource) *userRepository {
	return &userRepository{
		source: source,
	}
}

func (u *userRepository) Create(ctx context.Context, user *entity.UserCreate) (*entity.UserID, error) {
	userId, err := u.source.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("can't create user: %w", err)
	}

	return userId, nil
}

func (u *userRepository) GetById(ctx context.Context, id *entity.UserID) (*entity.User, error) {
	user, err := u.source.GetUserById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetByUsername: can't get user by id from db: %w", err)
	}

	return &entity.User{
		ID: &entity.UserID{
			Id: user.ID,
		},
		Username: user.Username,
	}, nil
}

func (u *userRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	user, err := u.source.GetUserByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("can't get user by id from db: %w", err)
	}

	return &entity.User{
		ID: &entity.UserID{
			Id: user.ID,
		},
		Username: user.Username,
	}, nil
}

func (u *userRepository) Update(ctx context.Context, id *entity.UserID, user *entity.UserCreate) (*entity.User, error) {
	_, err := u.source.GetUserById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("can't get user from db: %w", err)
	}

	dbUser, err := u.source.UpdateUser(ctx, id, user)
	if err != nil {
		return nil, fmt.Errorf("Update: can't to update user: %w", err)
	}

	return &entity.User{
		ID: &entity.UserID{
			Id: dbUser.ID,
		},
		Username: dbUser.Username,
	}, nil
}

func (u *userRepository) Delete(ctx context.Context, id *entity.UserID) error {
	_, err := u.source.GetUserById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return fmt.Errorf("can't get user from db: %w", err)
	}

	err = u.source.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("can't delete user from db: %w", err)
	}

	return nil
}

func (u *userRepository) LikeTrack(ctx context.Context, userId *entity.UserID, trackId *entity.MusicID) error {
	err := u.source.LikeTrack(ctx, userId, trackId)
	if err != nil {
		return fmt.Errorf("/db/user.LikeTrack: %w", err)
	}

	return nil
}

func (u *userRepository) DislikeTrack(ctx context.Context, userId *entity.UserID, trackId *entity.MusicID) error {
	err := u.source.DislikeTrack(ctx, userId, trackId)
	if err != nil {
		return fmt.Errorf("/db/user.DislikeTrack: %w", err)
	}

	return nil
}

func (u *userRepository) ShowLikedTracks(ctx context.Context, id *entity.UserID) ([]entity.MusicShow, error) {
	data, err := u.source.ShowLikedTracks(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("/db/user.ShowLikedTracks: %w", err)
	}

	return data, nil
}
