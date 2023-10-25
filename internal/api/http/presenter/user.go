package presenter

import (
	"fmt"
	"go-test-grpc-http/internal/api/http/view"
	"go-test-grpc-http/internal/entity"
)

type userPresenter struct {
}

func NewUserPresenter() *userPresenter {
	return &userPresenter{}
}

func (u *userPresenter) ToUserView(user *entity.User) *view.UserView {
	return &view.UserView{
		ID:    user.ID.String(),
		Name:  fmt.Sprintf("%s %s %s", user.LastName, user.FirstName, user.SecondName),
		Age:   user.Age,
		Email: user.Email,
		Phone: user.Phone,
	}
}
