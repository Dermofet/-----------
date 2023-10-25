package presenter

import (
	userv1 "go-test-grpc-http/internal/api/grpc/gen/servertemplate/user/v1"
	"go-test-grpc-http/internal/entity"
)

type UserPresenter interface {
	ToUser(user *userv1.UserDB) *entity.User
	FromUser(user *entity.User) *userv1.UserDB

	ToUserID(id string) *entity.UserID

	ToUserCreate(user *userv1.UserCreate) *entity.UserCreate
	FromUserCreate(user *entity.UserCreate) *userv1.UserCreate
}
