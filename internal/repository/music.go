package repository

import (
	"context"
	"music-backend-test/internal/db"
	"music-backend-test/internal/entity"
)

type musicRepository struct {
	source db.MusicSource
}

func NewMusicRepositiry(source db.MusicSource) *musicRepository {
	return &musicRepository{
		source: source,
	}
}

func (m *musicRepository) GetAll(ctx context.Context) ([]entity.MusicShow, error) {
	music, err := m.source.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return music, nil
}

func (m *musicRepository) Create(ctx context.Context, musicCreate *entity.MusicCreate) error {
	err := m.source.Create(ctx, musicCreate)
	if err != nil {
		return err
	}

	return nil
}

func (m *musicRepository) Update(ctx context.Context, musicUpdate *entity.MusicDB) error {
	err := m.source.Update(ctx, musicUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (m *musicRepository) Delete(ctx context.Context, id *entity.MusicID) error {
	err := m.source.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
