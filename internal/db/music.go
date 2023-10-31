package db

import (
	"context"
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

func (m *musicSource) GetAll(ctx context.Context) ([]*entity.Music, error) {
	var data []*entity.Music
	rows, err := m.db.Query("SELECT name FROM music")
	if err != nil {
		return nil, err
	}
	for i := 0; rows.Next(); i++ {
		var scanEntity entity.Music
		rows.Scan(&scanEntity.Name)
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
