package db

import (
	"context"
	"database/sql"
	"fmt"
	"music-backend-test/internal/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type musicSource struct {
	db *sqlx.DB
}

func NewMusicSource(source *source) *musicSource {
	return &musicSource{
		db: source.db,
	}
}

func (m *musicSource) GetAll(ctx context.Context) ([]*entity.MusicDB, error) {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	rows, err := m.db.QueryxContext(dbCtx, "SELECT * FROM music")
	if err != nil {
		return nil, fmt.Errorf("can't exec query: %w", err)
	}

	var data []*entity.MusicDB
	for i := 0; rows.Next(); i++ {
		var scanEntity entity.MusicDB
		err := rows.StructScan(&scanEntity)
		if err != nil {
			return nil, fmt.Errorf("can't scan music: %w", err)
		}
		data = append(data, &scanEntity)
	}
	return data, nil
}

func (m *musicSource) Get(ctx context.Context, musicId uuid.UUID) (*entity.MusicDB, error) {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	row := m.db.QueryRowxContext(dbCtx, "SELECT * FROM music WHERE id = $1", musicId)
	if row.Err() != nil {
		return nil, fmt.Errorf("can't exec query: %w", row.Err())
	}

	var data entity.MusicDB
	if err := row.StructScan(&data); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("can't scan music: %w", err)
	}

	return &data, nil
}

func (m *musicSource) GetAndSortByPopular(ctx context.Context) ([]*entity.MusicDB, error) {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	rows, err := m.db.QueryxContext(dbCtx,
		"SELECT m.id AS id, m.name AS name "+
			"FROM music m "+
			"LEFT JOIN user_music um ON um.music_id = m.id "+
			"GROUP BY m.id, m.name "+
			"ORDER BY COALESCE(COUNT(um.music_id), 0) DESC;",
	)
	if err != nil {
		return nil, fmt.Errorf("can't exec query: %w", err)
	}

	var data []*entity.MusicDB
	for i := 0; rows.Next(); i++ {
		var scanEntity entity.MusicDB
		err := rows.StructScan(&scanEntity)
		if err != nil {
			return nil, fmt.Errorf("can't scan music: %w", err)
		}
		data = append(data, &scanEntity)
	}

	return data, nil
}

func (m *musicSource) GetAllSortByTime(ctx context.Context) ([]*entity.MusicDB, error) {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	rows, err := m.db.QueryxContext(dbCtx, "SELECT * FROM music ORDER BY release_date")
	if err != nil {
		return nil, fmt.Errorf("can't exec query: %w", err)
	}

	var data []*entity.MusicDB
	for i := 0; rows.Next(); i++ {
		var scanEntity entity.MusicDB
		err := rows.StructScan(&scanEntity)
		if err != nil {
			return nil, fmt.Errorf("can't scan music: %w", err)
		}
		data = append(data, &scanEntity)
	}

	return data, nil
}

func (m *musicSource) Create(ctx context.Context, musicDb *entity.MusicDB) error {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	musicDb.Id = uuid.New()
	_, err := m.db.ExecContext(dbCtx, "INSERT INTO music (id, name, release_date, file_name, size, duration) VALUES ($1, $2, $3, $4, $5, $6)",
		musicDb.Id, musicDb.Name, musicDb.Release, musicDb.FileName, musicDb.Size, musicDb.Duration)
	if err != nil {
		return fmt.Errorf("can't exec query: %w", err)
	}

	return nil
}

func (m *musicSource) Update(ctx context.Context, musicDb *entity.MusicDB) error {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()
	if musicDb.FileName != "" {
		_, err := m.db.ExecContext(dbCtx, "UPDATE music SET name = $2, release_date = $3, file_name = $4, size = $5, duration = $6 WHERE id = $1",
			musicDb.Id, musicDb.Name, musicDb.Release, musicDb.FileName, musicDb.Size, musicDb.Duration)
		if err != nil {
			return fmt.Errorf("can't exec query: %w", err)
		}
	} else {
		_, err := m.db.ExecContext(dbCtx, "UPDATE music SET name = $2, release_date = $3 WHERE id = $1",
			musicDb.Id, musicDb.Name, musicDb.Release)
		if err != nil {
			return fmt.Errorf("can't exec query: %w", err)
		}
	}

	return nil
}

func (m *musicSource) Delete(ctx context.Context, id uuid.UUID) error {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()
	_, err := m.db.ExecContext(dbCtx, "DELETE FROM music WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("can't exec query: %w", err)
	}

	return nil
}
