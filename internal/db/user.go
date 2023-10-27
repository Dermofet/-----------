package db

import (
	"context"
	"database/sql"
	"fmt"
	"music-backend-test/internal/entity"

	"github.com/google/uuid"
)

func (s *source) CreateUser(ctx context.Context, user *entity.UserCreate) (*entity.UserID, error) {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	newUuid := uuid.New()

	row := s.db.QueryRowxContext(dbCtx, "INSERT INTO users (id, username, password) VALUES ($1, $2, $3)",
		newUuid, user.Username, user.Password)
	if row.Err() != nil {
		return nil, fmt.Errorf("can't exec query: %w", row.Err())
	}

	return &entity.UserID{
		Id: newUuid,
	}, nil
}

func (s *source) GetUserById(ctx context.Context, id *entity.UserID) (*entity.UserDB, error) {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	row := s.db.QueryRowxContext(dbCtx, "SELECT * FROM users WHERE id = $1", id.String())
	if row.Err() != nil {
		return nil, fmt.Errorf("can't exec query: %w", row.Err())
	}

	var userDB entity.UserDB
	if err := row.StructScan(&userDB); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("can't scan user: %w", err)
	}

	return &userDB, nil
}

func (s *source) GetUserByUsername(ctx context.Context, username string) (*entity.UserDB, error) {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	row := s.db.QueryRowxContext(dbCtx, "SELECT * FROM users WHERE username = $1", username)
	if row.Err() != nil {
		return nil, fmt.Errorf("can't exec query: %w", row.Err())
	}

	var userDB entity.UserDB
	if err := row.StructScan(&userDB); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("can't scan user: %w", err)
	}

	return &userDB, nil
}

func (s *source) UpdateUser(ctx context.Context, id *entity.UserID, user *entity.UserCreate) (*entity.UserDB, error) {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	row := s.db.QueryRowxContext(dbCtx, "UPDATE users SET username = $1, password = $2 WHERE id = $3",
		user.Username, user.Password, id.String())
	if row.Err() != nil {
		return nil, fmt.Errorf("can't exec query: %w", row.Err())
	}

	return &entity.UserDB{
		ID:       id.Id,
		Username: user.Username,
		Password: user.Password,
	}, nil
}

func (s *source) DeleteUser(ctx context.Context, id *entity.UserID) error {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	row := s.db.QueryRowxContext(dbCtx, "DELETE FROM users WHERE id = $1", id.String())
	if row.Err() != nil {
		return fmt.Errorf("can't exec query: %w", row.Err())
	}

	return nil
}
