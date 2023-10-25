package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-test-grpc-http/internal/api/http/presenter"
	_ "go-test-grpc-http/internal/api/http/view"
	"go-test-grpc-http/internal/entity"
	"go-test-grpc-http/internal/usecase"
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

// GetByEmailHandler godoc
// @Summary Получение пользователя по Email
// @Description Получение информации о пользователе по его уникальному идентификатору.
// @Tags Users
// @Accept json
// @Produce json
// @Param email path string true "Email пользователя"
// @Security JwtAuth
// @Success 200 {object} view.UserView "Данные пользователя"
// @Failure 400 "Некорректный запрос"
// @Failure 401 "Неавторизованный запрос"
// @Failure 404 "Пользователь не найден"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /users/email/{email} [get]
func (h *userHandlers) GetByEmailHandler(c *gin.Context) {
	ctx := context.Background()

	email := c.Param("email")
	user, err := h.interactor.GetByEmail(ctx, email)
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
// @Router /users/id/{id} [put]
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
// @Router /users/id/{id} [delete]
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
