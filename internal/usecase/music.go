package usecase

import (
	"context"
	"fmt"
	"music-backend-test/internal/entity"
	"music-backend-test/internal/repository"
)

type musicInteractor struct {
	repo repository.MusicRepository
}

func NewMusicInteractor(repo repository.MusicRepository) *musicInteractor {
	return &musicInteractor{
		repo: repo,
	}
}

func (m *musicInteractor) GetAll(ctx context.Context) ([]*entity.Music, error) {
	music, err := m.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("/repository/music.GetAll: %w", err)
	}

	return music, nil
}

func (m *musicInteractor) Get(ctx context.Context, musicId *entity.MusicID) (*entity.Music, error) {
	music, err := m.repo.Get(ctx, musicId)
	if err != nil {
		return nil, fmt.Errorf("/repository/music.Get: %w", err)
	}

	return music, nil
}

func (m *musicInteractor) GetAllSortByTime(ctx context.Context) ([]*entity.Music, error) {
	musics, err := m.repo.GetAllSortByTime(ctx)
	if err != nil {
		return nil, fmt.Errorf("/repository/music.GetAllSortByTime: %w", err)
	}
	return musics, nil
}

func (m musicInteractor) GetAndSortByPopular(ctx context.Context) ([]*entity.Music, error) {
	musics, err := m.repo.GetAndSortByPopular(ctx)
	if err != nil {
		return nil, fmt.Errorf("/repository/music.GetAndSortByPopular: %w", err)
	}
	return musics, nil
}

func (m *musicInteractor) Create(ctx context.Context, musicParse *entity.MusicParse) error {
	err := m.repo.Create(ctx, musicParse)
	if err != nil {
		return fmt.Errorf("/repository/music.Create: %w", err)
	}

	return nil
}

func (m *musicInteractor) Update(ctx context.Context, musicParse *entity.MusicParse) error {
	err := m.repo.Update(ctx, musicParse)
	if err != nil {
		return fmt.Errorf("/repository/music.Update: %w", err)
	}

	return nil
}

func (m *musicInteractor) Delete(ctx context.Context, id *entity.MusicID) error {
	err := m.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("/repository/music.Delete: %w", err)
	}

	return nil
}
