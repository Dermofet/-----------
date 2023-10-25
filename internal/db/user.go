package db

import (
	"context"
	"database/sql"
	"fmt"
	"go-test-grpc-http/internal/entity"

	"github.com/google/uuid"
)

func (s *source) CreateUser(ctx context.Context, user *entity.UserCreate) (*entity.UserID, error) {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	newUuid := uuid.New()

	row := s.db.QueryRowxContext(dbCtx, "INSERT INTO users (id, first_name, last_name, second_name, age, email, phone, password) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		newUuid, user.FirstName, user.LastName, user.SecondName, user.Age, user.Email, user.Phone, user.Password)
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

func (s *source) GetUserByEmail(ctx context.Context, email string) (*entity.UserDB, error) {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	row := s.db.QueryRowxContext(dbCtx, "SELECT * FROM users WHERE email = $1", email)
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

func (s *source) GetUserIdByEmail(ctx context.Context, email string) (*entity.UserID, error) {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	row := s.db.QueryRowxContext(dbCtx, "SELECT id FROM users WHERE email = $1", email)
	if row.Err() != nil {
		return nil, fmt.Errorf("can't exec query: %w", row.Err())
	}

	var userId_str string
	if err := row.Scan(&userId_str); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("can't scan user id: %w", err)
	}

	userId := uuid.MustParse(userId_str)

	return &entity.UserID{
		Id: userId,
	}, nil
}

func (s *source) UpdateUser(ctx context.Context, id *entity.UserID, user *entity.UserCreate) (*entity.UserDB, error) {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	_, err := s.GetUserById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("can't get user: %w", err)
	}

	row := s.db.QueryRowxContext(dbCtx, "UPDATE users SET first_name = $1, last_name = $2, second_name = $3, age = $4, email = $5, phone = $6, password = $7 WHERE id = $8",
		user.FirstName, user.LastName, user.SecondName, user.Age, user.Email, user.Phone, user.Password, id.String())
	if row.Err() != nil {
		return nil, fmt.Errorf("can't exec query: %w", row.Err())
	}

	return &entity.UserDB{
		ID:         id.Id,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		SecondName: user.SecondName,
		Age:        user.Age,
		Email:      user.Email,
		Phone:      user.Phone,
		Password:   user.Password,
	}, nil
}

func (s *source) DeleteUser(ctx context.Context, id *entity.UserID) error {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	_, err := s.GetUserById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return fmt.Errorf("can't get user: %w", err)
	}

	row := s.db.QueryRowxContext(dbCtx, "DELETE FROM users WHERE id = $1", id.String())
	if row.Err() != nil {
		return fmt.Errorf("can't exec query: %w", row.Err())
	}

	return nil
}
