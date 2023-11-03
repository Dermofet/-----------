package repository

import (
	"context"
	"fmt"
	"io"
	"music-backend-test/internal/db"
	"music-backend-test/internal/entity"
	"os"
)

type musicRepository struct {
	source db.MusicSource
}

func NewMusicRepository(source db.MusicSource) *musicRepository {
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

func (m *musicRepository) Get(ctx context.Context, musicId *entity.MusicID) (*entity.Music, error) {
	musicDB, err := m.source.Get(ctx, musicId)
	if err != nil {
		return nil, fmt.Errorf("/db/music.GetAll: %w", err)
	}

	music := &entity.Music{
		Id:   &entity.MusicID{Id: musicDB.Id},
		Name: musicDB.Name,
	}

	return music, nil
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

func (m *musicRepository) Create(ctx context.Context, musicParse *entity.MusicParse) error {
	musicCreate := &entity.MusicDB{
		Name:    musicParse.Name,
		Release: musicParse.Release,
	}

	download_file, err := os.Create("internal/storage/music_storage/" + musicParse.FileHeader.Filename)
	if err != nil {
		return fmt.Errorf("can't create file: %w", err)
	}
	defer download_file.Close()
	defer musicParse.File.Close()

	if _, err := io.Copy(download_file, musicParse.File); err != nil {
		return fmt.Errorf("can't copy file: %w", err)
	}
	musicCreate.FileName = musicParse.FileHeader.Filename

	err = m.source.Create(ctx, musicCreate)
	if err != nil {
		return fmt.Errorf("/db/music.Create: %w", err)
	}

	return nil
}

func (m *musicRepository) Update(ctx context.Context, musicParse *entity.MusicParse) error {
	musicUpdate := &entity.MusicDB{
		Id:      musicParse.Id,
		Name:    musicParse.Name,
		Release: musicParse.Release,
	}

	if musicParse.File != nil {
		music, err := m.source.Get(ctx, &entity.MusicID{Id: musicParse.Id})
		if err != nil {
			return fmt.Errorf("/db/music.Get: %w", err)
		}

		os.Remove("internal/storage/music_storage/" + music.FileName)
		download_file, err := os.Create("internal/storage/music_storage/" + musicParse.FileHeader.Filename)
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

func (m *musicRepository) Delete(ctx context.Context, id *entity.MusicID) error {
	music, err := m.source.Get(ctx, &entity.MusicID{Id: id.Id})
	if err != nil {
		return fmt.Errorf("/db/music.Get: %w", err)
	}
	os.Remove("internal/storage/music_storage/" + music.FileName)

	err = m.source.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("/db/music.Delete: %w", err)
	}

	return nil
}

//Выдача файла, тут кривость, тк сделано через http, но мб по сути этого можно найти аналог на gin
//Описание Header
//w.Header().Set("Content-Disposition", "attachment; filename="+file.Name)
//w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
//Отпавка файла
//http.ServeFile(w, r, "internal/storage/files_storage/"+file.Name)
