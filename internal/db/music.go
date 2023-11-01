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
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	rows, err := m.db.QueryxContext(dbCtx, "Select name from music")
	if err != nil {
		return nil, err
	}

	for i := 0; rows.Next(); i++ {
		var scanEntity entity.MusicShow
		err := rows.StructScan(&scanEntity)
		if err != nil {
			return nil, err
		}
		data = append(data, scanEntity)
	}

	return data, nil
}

func (m *musicSource) GetAllSortByTime(ctx context.Context) ([]entity.MusicShow, error) {
	var data []entity.MusicShow
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	rows, err := m.db.QueryxContext(dbCtx, "Select * from music order by release_date")
	if err != nil {
		return nil, err
	}

	for i := 0; rows.Next(); i++ {
		var scanEntity entity.MusicShow
		err := rows.StructScan(&scanEntity)
		if err != nil {
			return nil, err
		}
		data = append(data, scanEntity)
	}

	return data, nil
}

func (m *musicSource) Create(ctx context.Context, musicCreate *entity.MusicCreate) error {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	newuuid := uuid.New()
	row := m.db.QueryRowxContext(dbCtx, "Insert into music (id, name, release_date) values($1, $2, $3)", newuuid, musicCreate.Name, musicCreate.Release.Time)
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (m *musicSource) Update(ctx context.Context, musicUpdate *entity.MusicDB) error {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	row := m.db.QueryRowxContext(dbCtx, "Update music set name = $1, release_date = $2 where id=$3", musicUpdate.Name, musicUpdate.Release.Time, musicUpdate.Id.String())
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}

func (m *musicSource) Delete(ctx context.Context, id *entity.MusicID) error {
	dbCtx, dbCancel := context.WithTimeout(ctx, QueryTimeout)
	defer dbCancel()

	row := m.db.QueryRowxContext(dbCtx, "Delete from music where id = $1", id.String())
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}
