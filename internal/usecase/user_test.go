package usecase

import (
	"context"
	"fmt"
	"go-test-grpc-http/internal/entity"
	"go-test-grpc-http/internal/repository"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func Test_userInteractor_GetById(t *testing.T) {
	type fields struct {
		userRepository *repository.MockUserRepository
	}
	type args struct {
		ctx context.Context
		id  *entity.UserID
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.User
		setup   func(a args, f fields)
		wantErr bool
	}{
		{
			name: "success GetById usecase",
			args: args{
				ctx: context.Background(),
				id: &entity.UserID{
					Id: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				},
			},
			want: &entity.User{
				ID: &entity.UserID{
					Id: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				},
				FirstName:  "John",
				LastName:   "Doe",
				SecondName: "DoeD",
				Age:        30,
				Email:      "doe@example.com",
				Phone:      "+1111111111",
				Password:   "qwerty1234",
			},
			setup: func(a args, f fields) {
				user := &entity.User{
					ID:         a.id,
					FirstName:  "John",
					LastName:   "Doe",
					SecondName: "DoeD",
					Age:        30,
					Email:      "doe@example.com",
					Phone:      "+1111111111",
					Password:   "qwerty1234",
				}
				f.userRepository.EXPECT().GetById(a.ctx, a.id).Return(user, nil)
			},
			wantErr: false,
		},
		{
			name: "error GetById usecase",
			args: args{
				ctx: context.Background(),
				id:  nil,
			},
			want: nil,
			setup: func(a args, f fields) {
				f.userRepository.EXPECT().GetById(a.ctx, a.id).Return(nil, fmt.Errorf("can't get user from repository"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				userRepository: repository.NewMockUserRepository(ctrl),
			}
			u := &userInteractor{
				repo: f.userRepository,
			}

			tt.setup(tt.args, f)

			got, err := u.GetById(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("userInteractor.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userInteractor.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userInteractor_Create(t *testing.T) {
	type fields struct {
		repo *repository.MockUserRepository
	}
	type args struct {
		ctx  context.Context
		user *entity.UserCreate
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.UserID
		setup   func(a args, f fields)
		wantErr bool
	}{
		{
			name: "success Create usecase",
			args: args{
				ctx: context.Background(),
				user: &entity.UserCreate{
					FirstName:  "John",
					LastName:   "Doe",
					SecondName: "DoeD",
					Age:        30,
					Email:      "doe@example.com",
					Phone:      "+1111111111",
					Password:   "qwerty1234",
				},
			},
			want: &entity.UserID{
				Id: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			setup: func(a args, f fields) {
				f.repo.EXPECT().Create(a.ctx, a.user).Return(&entity.UserID{
					Id: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "error Create usecase",
			args: args{
				ctx: context.Background(),
				user: &entity.UserCreate{
					FirstName:  "John",
					LastName:   "Doe",
					SecondName: "DoeD",
					Age:        30,
					Email:      "doe@example.com",
					Phone:      "+1111111111",
					Password:   "qwerty1234",
				},
			},
			want: nil,
			setup: func(a args, f fields) {
				f.repo.EXPECT().Create(a.ctx, a.user).Return(nil, fmt.Errorf("can't create user in repository"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo: repository.NewMockUserRepository(ctrl),
			}
			u := &userInteractor{
				repo: f.repo,
			}

			tt.setup(tt.args, f)

			got, err := u.Create(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("userInteractor.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userInteractor.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userInteractor_Update(t *testing.T) {
	type fields struct {
		repo *repository.MockUserRepository
	}
	type args struct {
		ctx  context.Context
		id   *entity.UserID
		user *entity.UserCreate
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.User
		setup   func(a args, f fields)
		wantErr bool
	}{
		{
			name: "success Update usecase",
			args: args{
				ctx: context.Background(),
				id: &entity.UserID{
					Id: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				},
				user: &entity.UserCreate{
					FirstName:  "",
					LastName:   "",
					SecondName: "",
					Age:        31,
					Email:      "",
					Phone:      "",
					Password:   "",
				},
			},
			want: &entity.User{
				ID: &entity.UserID{
					Id: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				},
				FirstName:  "John",
				LastName:   "Doe",
				SecondName: "DoeD",
				Age:        31,
				Email:      "doe@example.com",
				Phone:      "+1111111111",
				Password:   "qwerty1234",
			},
			setup: func(a args, f fields) {
				f.repo.EXPECT().Update(a.ctx, a.id, a.user).Return(&entity.User{
					ID: &entity.UserID{
						Id: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
					},
					FirstName:  "John",
					LastName:   "Doe",
					SecondName: "DoeD",
					Age:        31,
					Email:      "doe@example.com",
					Phone:      "+1111111111",
					Password:   "qwerty1234",
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "error Update usecase",
			args: args{
				ctx: context.Background(),
				id: &entity.UserID{
					Id: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				},
				user: &entity.UserCreate{
					FirstName:  "",
					LastName:   "",
					SecondName: "",
					Age:        31,
					Email:      "",
					Phone:      "",
					Password:   "",
				},
			},
			want: nil,
			setup: func(a args, f fields) {
				f.repo.EXPECT().Update(a.ctx, a.id, a.user).Return(nil, fmt.Errorf("can't update user in repository"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo: repository.NewMockUserRepository(ctrl),
			}
			u := &userInteractor{
				repo: f.repo,
			}

			tt.setup(tt.args, f)

			got, err := u.Update(tt.args.ctx, tt.args.id, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("userInteractor.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userInteractor.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userInteractor_Delete(t *testing.T) {
	type fields struct {
		repo *repository.MockUserRepository
	}
	type args struct {
		ctx context.Context
		id  *entity.UserID
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.User
		setup   func(a args, f fields)
		wantErr bool
	}{
		{
			name: "success Delete usecase",
			args: args{
				ctx: context.Background(),
				id: &entity.UserID{
					Id: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				},
			},
			want: nil,
			setup: func(a args, f fields) {
				f.repo.EXPECT().Delete(a.ctx, a.id).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error Delete usecase",
			args: args{
				ctx: context.Background(),
				id: &entity.UserID{
					Id: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				},
			},
			want: nil,
			setup: func(a args, f fields) {
				f.repo.EXPECT().Delete(a.ctx, a.id).Return(fmt.Errorf("can't update user in repository"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo: repository.NewMockUserRepository(ctrl),
			}
			u := &userInteractor{
				repo: f.repo,
			}

			tt.setup(tt.args, f)

			if err := u.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("userInteractor.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
