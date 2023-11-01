package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"music-backend-test/internal/api/http/presenter"
	"music-backend-test/internal/entity"
	"music-backend-test/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type musicHandlers struct {
	interactor usecase.MusicInteractor
	presenter  presenter.MusicPresenter
}

func NewMusicHandlers(interactor usecase.MusicInteractor, presenter presenter.MusicPresenter) *musicHandlers {
	return &musicHandlers{
		interactor: interactor,
		presenter:  presenter,
	}
}

func (m *musicHandlers) GetAll(c *gin.Context) {
	ctx := context.Background()
	musics, err := m.interactor.GetAll(ctx)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("can't get user: %w", err))
	}
	c.JSON(http.StatusOK, m.presenter.ToListMusicView(musics)) //Проверить вывод
}

func (m *musicHandlers) Create(c *gin.Context) {
	ctx := context.Background()
	body, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("can't read body: %w", err))
		return
	}

	var music entity.MusicCreate
	err = json.Unmarshal(body, &music)
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("can't unmarshal body: %w", err))
		return
	}
	err = m.interactor.Create(ctx, &music)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("can't create music: %w", err))
	}
	c.JSON(http.StatusOK, nil)
}

func (m *musicHandlers) Update(c *gin.Context) {
	ctx := context.Background()
	body, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("can't read body: %w", err))
		return
	}

	var music entity.MusicDB
	music.Id, err = uuid.Parse(c.Param("id"))
	err = json.Unmarshal(body, &music)
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("can't unmarshal body: %w", err))
		return
	}
	err = m.interactor.Update(ctx, &music)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("can't update music: %w", err))
	}
	c.JSON(http.StatusOK, nil)
}

func (m *musicHandlers) Delete(c *gin.Context) {
	ctx := context.Background()

	var musicId entity.MusicID
	var err error
	musicId.Id, err = uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("can't unmarshal body: %w", err))
		return
	}
	err = m.interactor.Delete(ctx, &musicId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("can't get user: %w", err))
	}
	c.JSON(http.StatusOK, nil)
}
