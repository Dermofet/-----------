package presenter

import (
	"music-backend-test/internal/api/http/view"
	"music-backend-test/internal/entity"
)

//go:generate mockgen -source=./interfaces.go -destination=./presenter_mock.go -package=presenter

type UserPresenter interface {
	ToUserView(user *entity.User) *view.UserView
}

type TokenPresenter interface {
	ToTokenView(token *entity.Token) (*view.TokenView, error)
}
