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

func (m *musicInteractor) GetAll(ctx context.Context) ([]entity.MusicShow, error) {
	music, err := m.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("/repository/music.GetAll: %w", err)
	}

	return music, nil
}

func (m *musicInteractor) GetAllSortByTime(ctx context.Context) ([]entity.MusicShow, error) {
	music, err := m.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("/repository/music.GetAll: %w", err)
	}

	return music, nil
}

func (m *musicInteractor) Create(ctx context.Context, musicCreate *entity.MusicCreate) error {
	err := m.repo.Create(ctx, musicCreate)
	if err != nil {
		return fmt.Errorf("/repository/music.Create: %w", err)
	}

	return nil
}

func (m *musicInteractor) Update(ctx context.Context, musicUpdate *entity.MusicDB) error {
	err := m.repo.Update(ctx, musicUpdate)
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
