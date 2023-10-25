package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go-test-grpc-http/internal/api/http/presenter"
	_ "go-test-grpc-http/internal/api/http/view"
	"go-test-grpc-http/internal/entity"
	"go-test-grpc-http/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authHandlers struct {
	interactor usecase.UserInteractor
	presenter  presenter.TokenPresenter
}

func NewAuthHandlers(interactor usecase.UserInteractor, presenter presenter.TokenPresenter) *authHandlers {
	return &authHandlers{
		interactor: interactor,
		presenter:  presenter,
	}
}

// SignUp godoc
// @Summary Регистрация пользователя
// @Description Регистрация нового пользователя.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body entity.UserCreate true "Данные пользователя для регистрации"
// @Success 201 {object} view.TokenView "Токен авторизации"
// @Failure 400 "Некорректный запрос"
// @Failure 422 "Ошибка при обработке данных"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /auth/signup [post]
func (a *authHandlers) SignUp(c *gin.Context) {
	ctx := context.Background()

	data, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("can't sign up user: %v", err))
		return
	}

	var user *entity.UserCreate
	err = json.Unmarshal(data, &user)
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("can't sign up user: %v", err))
		return
	}

	userId, err := a.interactor.GetIdByEmail(ctx, user.Email)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("can't get user from db: %v", err))
		return
	}

	if userId != nil {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	userId, err = a.interactor.Create(ctx, user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("can't sign up user: %v", err))
		return
	}

	token, err := a.presenter.ToTokenView(entity.GenerateToken(userId))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("can't sign up user: %v", err))
		return
	}

	c.JSON(http.StatusOK, token)
}

// SignIn godoc
// @Summary Вход пользователя
// @Description Авторизация пользователя с использованием email и пароля.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body entity.UserSignIn true "Данные пользователя для входа"
// @Success 200 {object} view.TokenView "Токен авторизации"
// @Failure 400 "Некорректный запрос"
// @Failure 401 "Ошибка авторизации"
// @Failure 422 "Ошибка при обработке данных"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /auth/signin [post]
func (a *authHandlers) SignIn(c *gin.Context) {
	ctx := context.Background()

	data, err := c.GetRawData()
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("can't sign up user: %v", err))
		return
	}

	var user *entity.UserSignIn
	err = json.Unmarshal(data, &user)
	if err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("can't sign up user: %v", err))
		return
	}

	userID, err := a.interactor.GetIdByEmail(ctx, user.Email)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("can't get user from interactor: %v", err))
		return
	}

	if userID == nil {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	token, err := a.presenter.ToTokenView(entity.GenerateToken(userID))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("can't sign in user: %v", err))
		return
	}

	c.JSON(http.StatusOK, token)
}
