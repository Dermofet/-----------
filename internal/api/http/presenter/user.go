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

func (p *userPresenter) ToUserView(user *entity.User) *view.UserView {
	return &view.UserView{
		Id:       user.ID.String(),
		Username: user.Username,
		Role:     user.Role,
	}
}
