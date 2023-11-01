package db

import (
	"context"
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
		rows.StructScan(&scanEntity)
		data = append(data, &scanEntity)
	}
	return data, nil
}

func (m *musicSource) GetAndSortByPopular(ctx context.Context) ([]*entity.MusicDB, error) {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	rows, err := m.db.QueryxContext(dbCtx,
		"SELECT m.id AS id, m.name AS name\n"+
			"FROM music m\n"+
			"LEFT JOIN user_music um ON um.music_id = m.id\n"+
			"GROUP BY m.id, m.name\n"+
			"ORDER BY COALESCE(COUNT(um.music_id), 0) DESC;",
	)
	if err != nil {
		return nil, fmt.Errorf("can't exec query: %w", err)
	}

	var data []*entity.MusicDB
	for i := 0; rows.Next(); i++ {
		var scanEntity entity.MusicDB
		rows.StructScan(&scanEntity)
		fmt.Println(scanEntity)
		data = append(data, &scanEntity)
	}
	return data, nil
}

func (m *musicSource) Create(ctx context.Context, musicCreate *entity.MusicCreate) error {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	newuuid := uuid.New()
	row := m.db.QueryRowxContext(dbCtx, "INSERT INTO music (id, name) VALUES ($1, $2)", newuuid, musicCreate.Name)
	if row.Err() != nil {
		return row.Err()
	}
	return nil
}

func (m *musicSource) Update(ctx context.Context, musicUpdate *entity.MusicDB) error {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()
	row := m.db.QueryRowxContext(dbCtx, "UPDATE music SET name = $1 WHERE id = $2", musicUpdate.Name, musicUpdate.Id)
	if row.Err() != nil {
		return row.Err()
	}
	return nil
}

func (m *musicSource) Delete(ctx context.Context, id *entity.MusicID) error {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()
	row := m.db.QueryRowxContext(dbCtx, "DELETE FROM music WHERE id = $1", id.Id)
	if row.Err() != nil {
		return row.Err()
	}
	return nil
}
