package presenter

import (
	"music-backend-test/internal/api/http/view"
	"music-backend-test/internal/entity"
)

type Presenter interface {
	ToUserView(user *entity.UserDB) *view.UserView
	ToListUserView(users []*entity.UserDB) []*view.UserView
	ToMusicView(*entity.MusicDB) *view.MusicView
	ToListMusicView([]*entity.MusicDB) []*view.MusicView
	ToTokenView(token *entity.Token) (*view.TokenView, error)
}
