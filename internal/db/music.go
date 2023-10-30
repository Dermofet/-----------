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

func (m *musicSource) GetAll(ctx context.Context) ([]entity.MusicShow, error) {
	var data []entity.MusicShow
	rows, err := m.db.Query("Select name from music")
	if err != nil {
		return nil, err
	}
	for i := 0; rows.Next(); i++ {
		var scanEntity entity.MusicShow
		rows.Scan(&scanEntity.Name)
		data = append(data, scanEntity)
	}
	return data, nil
}

func (m *musicSource) Create(ctx context.Context, musicCreate *entity.MusicCreate) error {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	newuuid := uuid.New()
	row := m.db.QueryRowxContext(dbCtx, "Insert into music (id, name) values($1, $2)", newuuid, musicCreate.Name)
	if row.Err() != nil {
		return row.Err()
	}
	return nil
}

func (m *musicSource) Update(ctx context.Context, musicUpdate *entity.MusicDB) error {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()
	row := m.db.QueryRowxContext(dbCtx, "Update music set name = $1 where id=$2", musicUpdate.Name, musicUpdate.Id)
	if row.Err() != nil {
		return row.Err()
	}
	return nil
}

func (m *musicSource) Delete(ctx context.Context, id *entity.MusicID) error {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()
	row := m.db.QueryRowxContext(dbCtx, "Delete from music where id = $1", id.Id)
	if row.Err() != nil {
		return row.Err()
	}
	return nil
}
