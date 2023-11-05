package db

import (
	"context"
	"database/sql"
	"fmt"
	"music-backend-test/internal/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserSourсe struct {
	db *sqlx.DB
}

func NewUserSourсe(source *source) *UserSourсe {
	return &UserSourсe{
		db: source.db,
	}
}

func (u *UserSourсe) CreateUser(ctx context.Context, user *entity.UserCreate) (uuid.UUID, error) {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	newUuid := uuid.New()

	row := u.db.QueryRowxContext(dbCtx, "INSERT INTO users (id, username, password) VALUES ($1, $2, $3)",
		newUuid, user.Username, user.Password)
	if row.Err() != nil {
		return uuid.Nil, fmt.Errorf("can't exec query: %w", row.Err())
	}

	return newUuid, nil
}

func (u *UserSourсe) GetUserById(ctx context.Context, id uuid.UUID) (*entity.UserDB, error) {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	row := u.db.QueryRowxContext(dbCtx, "SELECT * FROM users WHERE id = $1", id.String())
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

func (u *UserSourсe) GetUserByUsername(ctx context.Context, username string) (*entity.UserDB, error) {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	row := u.db.QueryRowxContext(dbCtx, "SELECT * FROM users WHERE username = $1", username)
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

func (u *UserSourсe) UpdateUser(ctx context.Context, userDB *entity.UserDB, user *entity.UserCreate) (*entity.UserDB, error) {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	row := u.db.QueryRowxContext(dbCtx, "UPDATE users SET username = $1, password = $2 WHERE id = $3",
		user.Username, user.Password, userDB.ID.String())
	if row.Err() != nil {
		return nil, fmt.Errorf("can't exec query: %w", row.Err())
	}

	userDB.Username = user.Username
	userDB.Password = user.Password
	return userDB, nil
}

func (u *UserSourсe) DeleteUser(ctx context.Context, id uuid.UUID) error {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	row := u.db.QueryRowxContext(dbCtx, "DELETE FROM users WHERE id = $1", id.String())
	if row.Err() != nil {
		return fmt.Errorf("can't exec query: %w", row.Err())
	}

	return nil
}

func (u *UserSourсe) LikeTrack(ctx context.Context, userId uuid.UUID, trackId uuid.UUID) error {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	row := u.db.QueryRowxContext(dbCtx, "Insert into user_music (user_id, music_id) values ($1, $2)", userId.String(), trackId.String())
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (u *UserSourсe) DislikeTrack(ctx context.Context, userId uuid.UUID, trackId uuid.UUID) error {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	row := u.db.QueryRowxContext(dbCtx, "Delete from user_music where user_id = $1 AND music_id = $2", userId.String(), trackId.String())
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (u *UserSourсe) ShowLikedTracks(ctx context.Context, id uuid.UUID) ([]*entity.MusicDB, error) {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	rows, err := u.db.QueryxContext(
		dbCtx,
		"SELECT music.* FROM music JOIN user_music ON music.id = user_music.music_id WHERE user_music.user_id = $1",
		id.String(),
	)
	if err != nil {
		return nil, err
	}

	var data []*entity.MusicDB
	for i := 0; rows.Next(); i++ {
		var scanEntity entity.MusicDB
		err := rows.StructScan(&scanEntity)
		if err != nil {
			return nil, err
		}
		fmt.Println(scanEntity)
		data = append(data, &scanEntity)
	}

	return data, nil
}
