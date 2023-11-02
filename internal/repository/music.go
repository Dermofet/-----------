package repository

import (
	"context"
	"fmt"
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

func (m *musicRepository) GetAll(ctx context.Context) ([]*entity.Music, error) {
	musicsDB, err := m.source.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("/db/music.GetAll: %w", err)
	}

	musics := make([]*entity.Music, len(musicsDB))
	for i, musicDB := range musicsDB {
		musics[i] = &entity.Music{
			Id:   &entity.MusicID{Id: musicDB.Id},
			Name: musicDB.Name,
		}
	}
	return musics, nil
}

func (m *musicRepository) GetAndSortByPopular(ctx context.Context) ([]*entity.Music, error) {
	musicsDB, err := m.source.GetAndSortByPopular(ctx)
	if err != nil {
		return nil, fmt.Errorf("/db/music.GetAndSortByPopular: %w", err)
	}

	musics := make([]*entity.Music, len(musicsDB))
	for i, musicDB := range musicsDB {
		musics[i] = &entity.Music{
			Id:   &entity.MusicID{Id: musicDB.Id},
			Name: musicDB.Name,
		}
	}
	return musics, nil
}

func (m *musicRepository) GetAllSortByTime(ctx context.Context) ([]*entity.Music, error) {
	musicsDB, err := m.source.GetAllSortByTime(ctx)
	if err != nil {
		return nil, fmt.Errorf("/db/music.GetAllSortByTime: %w", err)
	}

	musics := make([]*entity.Music, len(musicsDB))
	for i, musicDB := range musicsDB {
		musics[i] = &entity.Music{
			Id:   &entity.MusicID{Id: musicDB.Id},
			Name: musicDB.Name,
		}
	}
	return musics, nil
}

func (m *musicRepository) Create(ctx context.Context, musicCreate *entity.MusicCreate) error {
	err := m.source.Create(ctx, musicCreate)
	if err != nil {
		return fmt.Errorf("/db/music.Create: %w", err)
	}

	return nil
}

func (m *musicRepository) Update(ctx context.Context, musicUpdate *entity.MusicDB) error {
	err := m.source.Update(ctx, musicUpdate)
	if err != nil {
		return fmt.Errorf("/db/music.Update: %w", err)
	}

	return nil
}

func (m *musicRepository) Delete(ctx context.Context, id *entity.MusicID) error {
	err := m.source.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("/db/music.Delete: %w", err)
	}

	return nil
}
