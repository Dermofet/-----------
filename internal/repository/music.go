package repository

import (
	"context"
	"fmt"
	"io"
	"music-backend-test/internal/db"
	"music-backend-test/internal/entity"
	"music-backend-test/internal/utils"
	"os"

	"github.com/google/uuid"
)

type musicRepository struct {
	source db.MusicSource
}

func NewMusicRepository(source db.MusicSource) *musicRepository {
	return &musicRepository{
		source: source,
	}
}

func (m *musicRepository) GetAll(ctx context.Context) ([]*entity.MusicDB, error) {
	musicsDB, err := m.source.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("/db/music.GetAll: %w", err)
	}

	// musics := make([]*entity.MusicDB, len(musicsDB))
	// for i, musicDB := range musicsDB {
	// 	musics[i] = &entity.Music{
	// 		Id:   musicDB.Id,
	// 		Name: musicDB.Name,
	// 	}
	// }
	return musicsDB, nil
}

func (m *musicRepository) Get(ctx context.Context, musicId uuid.UUID) (*entity.MusicDB, error) {
	musicDB, err := m.source.Get(ctx, musicId)
	if err != nil {
		return nil, fmt.Errorf("/db/music.Get: %w", err)
	}

	// music := &entity.Music{
	// 	Id:   musicDB.Id,
	// 	Name: musicDB.Name,
	// }

	return musicDB, nil
}

func (m *musicRepository) GetAndSortByPopular(ctx context.Context) ([]*entity.MusicDB, error) {
	musicsDB, err := m.source.GetAndSortByPopular(ctx)
	if err != nil {
		return nil, fmt.Errorf("/db/music.GetAndSortByPopular: %w", err)
	}

	// musics := make([]*entity.MusicDB, len(musicsDB))
	// for i, musicDB := range musicsDB {
	// 	musics[i] = &entity.Music{
	// 		Id:   musicDB.Id,
	// 		Name: musicDB.Name,
	// 	}
	// }
	return musicsDB, nil
}

func (m *musicRepository) GetAllSortByTime(ctx context.Context) ([]*entity.MusicDB, error) {
	musicsDB, err := m.source.GetAllSortByTime(ctx)
	if err != nil {
		return nil, fmt.Errorf("/db/music.GetAllSortByTime: %w", err)
	}

	// musics := make([]*entity.MusicDB, len(musicsDB))
	// for i, musicDB := range musicsDB {
	// 	musics[i] = &entity.Music{
	// 		Id:   musicDB.Id,
	// 		Name: musicDB.Name,
	// 	}
	// }
	return musicsDB, nil
}

func (m *musicRepository) Create(ctx context.Context, musicParse *entity.MusicParse, fileType utils.FileType) error {
	musicCreate := &entity.MusicDB{
		Name:     musicParse.Name,
		Release:  musicParse.Release,
		FileName: musicParse.FileHeader.Filename,
	}

	var err error
	musicCreate.Size = uint64(musicParse.FileHeader.Size)

	musicCreate.Duration, err = utils.GetAudioDuration(fileType, musicCreate.FilePath())
	if err != nil {
		return fmt.Errorf("/utils.GetAudioDuration: %w", err)
	}

	// скачиваем файл
	download_file, err := os.Create(musicCreate.FilePath())
	if err != nil {
		return fmt.Errorf("can't create file: %w", err)
	}
	defer download_file.Close()
	defer musicParse.File.Close()

	if _, err := io.Copy(download_file, musicParse.File); err != nil {
		return fmt.Errorf("can't copy file: %w", err)
	}

	err = m.source.Create(ctx, musicCreate)
	if err != nil {
		return fmt.Errorf("/db/music.Create: %w", err)
	}

	return nil
}

func (m *musicRepository) Update(ctx context.Context, id uuid.UUID, musicParse *entity.MusicParse, fileType utils.FileType) error {
	musicUpdate := &entity.MusicDB{
		Id:      id,
		Name:    musicParse.Name,
		Release: musicParse.Release,
	}

	if musicParse.FileHeader != nil {
		music, err := m.source.Get(ctx, id)
		if err != nil {
			return fmt.Errorf("/db/music.Get: %w", err)
		}

		musicUpdate.Size = uint64(musicParse.FileHeader.Size)

		musicUpdate.Duration, err = utils.GetAudioDuration(fileType, musicUpdate.FilePath())
		if err != nil {
			return fmt.Errorf("/utils.GetAudioDuration: %w", err)
		}

		os.Remove(music.FilePath())
		download_file, err := os.Create(music.FilePath())
		if err != nil {
			return fmt.Errorf("can't create file: %w", err)
		}
		defer download_file.Close()
		defer musicParse.File.Close()

		if _, err := io.Copy(download_file, musicParse.File); err != nil {
			return fmt.Errorf("can't copy file: %w", err)
		}
		musicUpdate.FileName = musicParse.FileHeader.Filename
	}

	err := m.source.Update(ctx, musicUpdate)
	if err != nil {
		return fmt.Errorf("/db/music.Update: %w", err)
	}

	return nil
}

func (m *musicRepository) Delete(ctx context.Context, id uuid.UUID) error {
	music, err := m.source.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("/db/music.Get: %w", err)
	}
	os.Remove(music.FilePath())

	err = m.source.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("/db/music.Delete: %w", err)
	}

	return nil
}
