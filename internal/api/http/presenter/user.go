package presenter

import (
	"music-backend-test/internal/api/http/view"
	"music-backend-test/internal/entity"
)

type userPresenter struct {
}

func NewUserPresenter() *userPresenter {
	return &userPresenter{}
}

func (u *userPresenter) ToUserView(user *entity.User) *view.UserView {
	return &view.UserView{
		ID:       user.ID.String(),
		Username: user.Username,
	}
}
