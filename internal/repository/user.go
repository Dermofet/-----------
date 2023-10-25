package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-test-grpc-http/internal/db"
	"go-test-grpc-http/internal/entity"
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
		return nil, fmt.Errorf("can't get user by id from db: %w", err)
	}

	return &entity.User{
		ID: &entity.UserID{
			Id: user.ID,
		},
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		LastName:   user.LastName,
		Password:   user.Password,
		Age:        user.Age,
		Email:      user.Email,
		Phone:      user.Phone,
	}, nil
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := u.source.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("can't get user by email from db: %w", err)
	}

	return &entity.User{
		ID: &entity.UserID{
			Id: user.ID,
		},
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		LastName:   user.LastName,
		Password:   user.Password,
		Age:        user.Age,
		Email:      user.Email,
		Phone:      user.Phone,
	}, nil
}

func (u *userRepository) GetIdByEmail(ctx context.Context, email string) (*entity.UserID, error) {
	id, err := u.source.GetUserIdByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("can't get user id by email from db: %w", err)
	}

	return id, nil
}

func (u *userRepository) Update(ctx context.Context, id *entity.UserID, user *entity.UserCreate) (*entity.User, error) {
	dbUser, err := u.source.UpdateUser(ctx, id, user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("can't to update user: %w", err)
	}

	return &entity.User{
		ID: &entity.UserID{
			Id: dbUser.ID,
		},
		FirstName:  dbUser.FirstName,
		SecondName: dbUser.SecondName,
		LastName:   dbUser.LastName,
		Password:   dbUser.Password,
		Age:        dbUser.Age,
		Email:      dbUser.Email,
		Phone:      dbUser.Phone,
	}, nil
}

func (u *userRepository) Delete(ctx context.Context, id *entity.UserID) error {
	err := u.source.DeleteUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return fmt.Errorf("can't delete user from db: %w", err)
	}

	return nil
}
