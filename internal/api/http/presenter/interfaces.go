package presenter

import (
	"music-backend-test/internal/api/http/view"
	"music-backend-test/internal/entity"
)

type UserPresenter interface {
	ToUserView(user *entity.User) *view.UserView
}

type TokenPresenter interface {
	ToTokenView(token *entity.Token) (*view.TokenView, error)
}

type MusicPresenter interface {
	ToMusicView(*entity.Music) *view.MusicView
	ToListMusicView([]*entity.Music) []*view.MusicView
}
