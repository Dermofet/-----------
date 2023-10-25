package presenter

import (
	userv1 "go-test-grpc-http/internal/api/grpc/gen/servertemplate/user/v1"
	"go-test-grpc-http/internal/entity"

	"github.com/google/uuid"
)

type userPresenter struct {
}

func NewUserPresenter() *userPresenter {
	return &userPresenter{}
}

// TODO почекать выходные данные, чтобы не отправлять лишнюю инфу о юзерах

func (u *userPresenter) FromUser(user *entity.User) *userv1.UserDB {
	return &userv1.UserDB{
		Id:         user.ID.String(),
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		LastName:   user.LastName,
		Age:        int32(user.Age),
		Password:   user.Password,
		Email:      user.Email,
		Phone:      user.Phone,
	}
}

func (u *userPresenter) ToUser(user *userv1.UserDB) *entity.User {
	return &entity.User{
		ID:         u.ToUserID(user.GetId()),
		FirstName:  user.GetFirstName(),
		SecondName: user.GetSecondName(),
		LastName:   user.GetLastName(),
		Age:        int(user.GetAge()),
		Password:   user.GetPassword(),
		Email:      user.GetEmail(),
		Phone:      user.GetPhone(),
	}
}

func (u *userPresenter) ToUserID(id string) *entity.UserID {
	newId, err := uuid.Parse(id)
	if err != nil {
		return nil
	}
	return &entity.UserID{
		Id: newId,
	}
}

func (u *userPresenter) ToUserCreate(user *userv1.UserCreate) *entity.UserCreate {
	return &entity.UserCreate{
		FirstName:  user.GetFirstName(),
		SecondName: user.GetSecondName(),
		LastName:   user.GetLastName(),
		Age:        int(user.GetAge()),
		Password:   user.GetPassword(),
		Email:      user.GetEmail(),
		Phone:      user.GetPhone(),
	}
}

func (u *userPresenter) FromUserCreate(user *entity.UserCreate) *userv1.UserCreate {
	return &userv1.UserCreate{
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		LastName:   user.LastName,
		Age:        int32(user.Age),
		Password:   user.Password,
		Email:      user.Email,
		Phone:      user.Phone,
	}
}
