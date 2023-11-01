package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"music-backend-test/internal/api/http/presenter"
	"music-backend-test/internal/entity"
	"music-backend-test/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userHandlers struct {
	interactor usecase.UserInteractor
	presenter  presenter.UserPresenter
}

func NewUserHandlers(interactor usecase.UserInteractor, presenter presenter.UserPresenter) *userHandlers {
	return &userHandlers{
		interactor: interactor,
		presenter:  presenter,
	}
}

// GetMeHandler godoc
// @Summary Получение пользователя по JWT токену
// @Description Получение пользователя по его уникальному идентификатору из JWT токена
// @Tags Users
// @Accept json
// @Produce plain
// @Security JwtAuth
// @Success 200 {object} view.UserView "Данные пользователя"
// @Failure 400 "Некорректный запрос"
// @Failure 401 "Неавторизованный запрос"
// @Failure 404 "Пользователь не найден"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /users/me [get]
func (h *userHandlers) GetMeHandler(c *gin.Context) {
	ctx := context.Background()

	id, exists := c.Get("user-id")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := h.interactor.GetById(ctx, id.(*entity.UserID))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("can't get user: %w", err))
		return
	}

	if user == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, h.presenter.ToUserView(user))
}

// UpdateMeHandler godoc
// @Summary Обновление пользователя по JWT токену
// @Description Обновление информации о пользователе по его уникальному идентификатору из JWT токена
// @Tags Users
// @Accept json
// @Produce json
// @Param request body entity.UserCreate true "Данные пользователя для обновления"
// @Security JwtAuth
// @Success 200 {object} view.UserView "Обновленные данные пользователя"
// @Failure 400 "Некорректный запрос"
// @Failure 401 "Неавторизованный запрос"
// @Failure 404 "Пользователь не найден"
// @Failure 422 "Ошибка при обработке данных"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /users/me [put]
func (h *userHandlers) UpdateMeHandler(c *gin.Context) {
	ctx := context.Background()

	id, exists := c.Get("user-id")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	body, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("can't read body: %w", err))
		return
	}

	var user entity.UserCreate
	err = json.Unmarshal(body, &user)
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("can't unmarshal body: %w", err))
		return
	}

	dbUser, err := h.interactor.Update(ctx, id.(*entity.UserID), &user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("can't update user: %w", err))
		return
	}

	if dbUser == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, h.presenter.ToUserView(dbUser))
}

// DeleteMeHandler godoc
// @Summary Удаление пользователя по JWT токену
// @Description Удаление пользователя по его уникальному идентификатору из JWT токена.
// @Tags Users
// @Accept json
// @Produce plain
// @Security JwtAuth
// @Success 204 "Пользователь успешно удален"
// @Failure 400 "Некорректный запрос"
// @Failure 401 "Неавторизованный запрос"
// @Failure 404 "Пользователь не найден"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /users/me [delete]
func (h *userHandlers) DeleteMeHandler(c *gin.Context) {
	ctx := context.Background()

	id, exists := c.Get("user-id")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err := h.interactor.Delete(ctx, id.(*entity.UserID))
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("can't delete user: %w", err))
		return
	}

	c.Status(http.StatusNoContent)
}

// GetByIdHandler godoc
// @Summary Получение пользователя по ID
// @Description Получение информации о пользователе по его уникальному идентификатору.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "Уникальный идентификатор пользователя (UUID)"
// @Security JwtAuth
// @Success 200 {object} view.UserView "Данные пользователя"
// @Failure 400 "Некорректный запрос"
// @Failure 401 "Неавторизованный запрос"
// @Failure 404 "Пользователь не найден"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /users/id/{id} [get]
func (h *userHandlers) GetByIdHandler(c *gin.Context) {
	ctx := context.Background()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("invalid id: %w", err))
		return
	}

	user, err := h.interactor.GetById(ctx, &entity.UserID{Id: id})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("can't get user: %w", err))
		return
	}

	if user == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, h.presenter.ToUserView(user))
}

// GetByUsernameHandler godoc
// @Summary Получение пользователя по Username
// @Description Получение информации о пользователе по его уникальному идентификатору.
// @Tags Users
// @Accept json
// @Produce json
// @Param username path string true "Username пользователя"
// @Security JwtAuth
// @Success 200 {object} view.UserView "Данные пользователя"
// @Failure 400 "Некорректный запрос"
// @Failure 401 "Неавторизованный запрос"
// @Failure 404 "Пользователь не найден"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /users/username/{username} [get]
func (h *userHandlers) GetByUsernameHandler(c *gin.Context) {
	ctx := context.Background()

	username := c.Param("username")
	user, err := h.interactor.GetByUsername(ctx, username)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("can't get user: %w", err))
		return
	}

	if user == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, h.presenter.ToUserView(user))
}

// UpdateHandler godoc
// @Summary Обновление пользователя по ID
// @Description Обновление информации о пользователе по его уникальному идентификатору.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "Уникальный идентификатор пользователя (UUID)"
// @Param request body entity.UserCreate true "Данные пользователя для обновления"
// @Security JwtAuth
// @Success 200 {object} view.UserView "Обновленные данные пользователя"
// @Failure 400 "Некорректный запрос"
// @Failure 401 "Неавторизованный запрос"
// @Failure 404 "Пользователь не найден"
// @Failure 422 "Ошибка при обработке данных"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /users/{id} [put]
func (h *userHandlers) UpdateHandler(c *gin.Context) {
	ctx := context.Background()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	body, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("can't read body: %w", err))
		return
	}

	var user entity.UserCreate
	err = json.Unmarshal(body, &user)
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("can't unmarshal body: %w", err))
		return
	}

	dbUser, err := h.interactor.Update(ctx, &entity.UserID{Id: id}, &user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("can't update user: %w", err))
		return
	}

	c.JSON(http.StatusOK, h.presenter.ToUserView(dbUser))
}

// DeleteHandler godoc
// @Summary Удаление пользователя по ID
// @Description Удаление пользователя по его уникальному идентификатору.
// @Tags Users
// @Accept json
// @Produce plain
// @Param id path string true "Уникальный идентификатор пользователя (UUID)"
// @Security JwtAuth
// @Success 204 "Пользователь успешно удален"
// @Failure 400 "Некорректный запрос"
// @Failure 401 "Неавторизованный запрос"
// @Failure 404 "Пользователь не найден"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /users/{id} [delete]
func (h *userHandlers) DeleteHandler(c *gin.Context) {
	ctx := context.Background()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	err = h.interactor.Delete(ctx, &entity.UserID{Id: id})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("can't delete user: %w", err))
		return
	}

	c.Status(http.StatusNoContent)
}

// LikeTrackHandler godoc
// @Summary Добавить трека в список понравившихся
// @Description Добавление трека в список понравившихся
// @Tags Users
// @Accept json
// @Produce plain
// @Security JwtAuth
// @Param id path string true "Уникальный идентификатор трека (UUID)"
// @Success 201 "Трек успешно добавлен в список понравившихся"
// @Failure 400 "Некорректный запрос"
// @Failure 401 "Неавторизованный запрос"
// @Failure 404 "Трек не найден"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /users/add-track/{id} [post]
func (h *userHandlers) LikeTrack(c *gin.Context) {
	ctx := context.Background()

	id, exists := c.Get("user-id")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	body, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("can't read body: %w", err))
		return
	}

	var trackId *entity.MusicID
	err = json.Unmarshal(body, &trackId)
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("can't unmarshal body: %w", err))
		return
	}

	err = h.interactor.LikeTrack(ctx, userId, trackId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("/usecase/user.LikeTrack: %w", err))
		return
	}

	c.JSON(http.StatusOK, nil)
}

// RemoveTrackHandler godoc
// @Summary Удалить трек из списка понравившихся
// @Description Удаление трека из списка понравившихся
// @Tags Users
// @Accept json
// @Produce plain
// @Param id path string true "Уникальный идентификатор трека (UUID)"
// @Security JwtAuth
// @Success 204 "Трек успешно удален из списка понравившихся"
// @Failure 400 "Некорректный запрос"
// @Failure 401 "Неавторизованный запрос"
// @Failure 404 "Трек не найден"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /users/remove-track/{id} [delete]
func (h *userHandlers) DislikeTrack(c *gin.Context) {
	ctx := context.Background()

	id, exists := c.Get("user-id")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	body, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("can't read body: %w", err))
		return
	}

	var trackId *entity.MusicID
	err = json.Unmarshal(body, &trackId)
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("can't unmarshal body: %w", err))
		return
	}

	err = h.interactor.DislikeTrack(ctx, userId, trackId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("/usecase/user.DislikeTrack: %w", err))
		return
	}

	c.JSON(http.StatusOK, nil)
}

// ShowLikedTracksHandler godoc
// @Summary Показать понравившиеся треки
// @Description Получение списка понравившихся треков
// @Tags Users
// @Accept json
// @Produce plain
// @Security JwtAuth
// @Success 200 {object} view.ListMusicView
// @Failure 400 "Некорректный запрос"
// @Failure 401 "Неавторизованный запрос"
// @Failure 404 "Трек не найден"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /users/get-liked-tracks [get]
func (h *userHandlers) ShowLikedTracks(c *gin.Context) {
	ctx := context.Background()

	id, exists := c.Get("user-id")
	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	data, err := h.interactor.ShowLikedTracks(ctx, id.(*entity.UserID))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("/usecase/user.ShowLikedTracks: %w", err))
		return
	}

	c.JSON(http.StatusOK, data)
}
