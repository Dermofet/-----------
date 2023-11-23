package usecase

import (
	"context"
	"fmt"
	"music-backend-test/internal/entity"
	"music-backend-test/internal/repository"
	"music-backend-test/internal/usecase"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

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
		want    uuid.UUID
		setup   func(a args, f fields)
		wantErr bool
	}{
		{
			name: "success Create usecase",
			args: args{
				ctx: context.Background(),
				user: &entity.UserCreate{
					Username: "John",
					Password: "qwerty1234",
				},
			},
			want: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			setup: func(a args, f fields) {
				f.repo.EXPECT().Create(a.ctx, a.user).Return(uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"), nil)
			},
			wantErr: false,
		},
		{
			name: "error Create usecase",
			args: args{
				ctx: context.Background(),
				user: &entity.UserCreate{
					Username: "John",
					Password: "qwerty1234",
				},
			},
			want: uuid.Nil,
			setup: func(a args, f fields) {
				f.repo.EXPECT().Create(a.ctx, a.user).Return(uuid.Nil, fmt.Errorf("can't create user in repository"))
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
			u := usecase.NewUserInteractor(f.repo)
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

func Test_userInteractor_GetById(t *testing.T) {
	type fields struct {
		repo *repository.MockUserRepository
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.UserDB
		setup   func(a args, f fields)
		wantErr bool
	}{
		{
			name: "success GetById usecase",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			want: &entity.UserDB{
				ID:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				Username: "John",
			},
			setup: func(a args, f fields) {
				user := &entity.UserDB{
					ID:       a.id,
					Username: "John",
				}
				f.repo.EXPECT().GetById(a.ctx, a.id).Return(user, nil)
			},
			wantErr: false,
		},
		{
			name: "error GetById usecase",
			args: args{
				ctx: context.Background(),
				id:  uuid.Nil,
			},
			want: nil,
			setup: func(a args, f fields) {
				f.repo.EXPECT().GetById(a.ctx, a.id).Return(nil, fmt.Errorf("can't get user from repository"))
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
			u := usecase.NewUserInteractor(f.repo)

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

func Test_userInteractor_GetByUsername(t *testing.T) {
	type fields struct {
		repo *repository.MockUserRepository
	}
	type args struct {
		ctx      context.Context
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.UserDB
		setup   func(a args, f fields)
		wantErr bool
	}{
		{
			name: "success GetByUsername usecase",
			args: args{
				ctx:      context.Background(),
				username: "John",
			},
			want: &entity.UserDB{
				ID:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				Username: "John",
			},
			setup: func(a args, f fields) {
				user := &entity.UserDB{
					ID:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
					Username: "John",
				}
				f.repo.EXPECT().GetByUsername(a.ctx, a.username).Return(user, nil)
			},
			wantErr: false,
		},
		{
			name: "error GetByUsername usecase",
			args: args{
				username: "",
			},
			want: nil,
			setup: func(a args, f fields) {
				f.repo.EXPECT().GetByUsername(a.ctx, a.username).Return(nil, fmt.Errorf("can't get user from repository"))
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
			u := usecase.NewUserInteractor(f.repo)

			tt.setup(tt.args, f)

			got, err := u.GetByUsername(tt.args.ctx, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("userInteractor.GetByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userInteractor.GetByUsername() = %v, want %v", got, tt.want)
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
		id   uuid.UUID
		user *entity.UserCreate
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.UserDB
		setup   func(a args, f fields)
		wantErr bool
	}{
		{
			name: "success Update usecase",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				user: &entity.UserCreate{
					Username: "Paul",
					Password: "",
				},
			},
			want: &entity.UserDB{
				ID:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				Username: "Paul",
			},
			setup: func(a args, f fields) {
				f.repo.EXPECT().Update(a.ctx, a.id, a.user).Return(&entity.UserDB{
					ID:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
					Username: "Paul",
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "error Update usecase",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				user: &entity.UserCreate{
					Username: "Paul",
					Password: "",
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
			u := usecase.NewUserInteractor(f.repo)
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
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.UserDB
		setup   func(a args, f fields)
		wantErr bool
	}{
		{
			name: "success Delete usecase",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
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
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
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
			u := usecase.NewUserInteractor(f.repo)
			tt.setup(tt.args, f)

			if err := u.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("userInteractor.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userInteractor_LikeTrack(t *testing.T) {
	type fields struct {
		repo *repository.MockUserRepository
	}
	type args struct {
		ctx     context.Context
		userId  uuid.UUID
		trackId uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		setup   func(a args, f fields)
		want    error
		wantErr bool
	}{
		{
			name: "success: LikeTrack userInteractor",
			args: args{
				ctx:     context.Background(),
				userId:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				trackId: uuid.MustParse("5b60e78f-b465-4cd6-b5d3-15b188f47a6a"),
			},
			setup: func(a args, f fields) {
				f.repo.EXPECT().LikeTrack(a.ctx, a.userId, a.trackId).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error: LikeTrack userInteractor",
			args: args{
				ctx:     context.Background(),
				userId:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				trackId: uuid.MustParse("5b60e78f-b465-4cd6-b5d3-15b188f47a6a"),
			},
			setup: func(a args, f fields) {
				f.repo.EXPECT().LikeTrack(a.ctx, a.userId, a.trackId).Return(fmt.Errorf("can't like track in repository"))
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
			u := usecase.NewUserInteractor(f.repo)

			tt.setup(tt.args, f)

			err := u.LikeTrack(tt.args.ctx, tt.args.userId, tt.args.trackId)
			if (err != nil) != tt.wantErr {
				t.Errorf("userInteractor.LikeTrack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userInteractor_DislikeTrack(t *testing.T) {
	type fields struct {
		repo *repository.MockUserRepository
	}
	type args struct {
		ctx     context.Context
		userId  uuid.UUID
		trackId uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    error
		setup   func(a args, f fields)
		wantErr bool
	}{
		{
			name: "success: DislikeTrack userInteractor",
			args: args{
				ctx:     context.Background(),
				userId:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				trackId: uuid.MustParse("5b60e78f-b465-4cd6-b5d3-15b188f47a6a"),
			},
			want: nil,
			setup: func(a args, f fields) {
				f.repo.EXPECT().DislikeTrack(a.ctx, a.userId, a.trackId).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "error: DislikeTrack userInteractor",
			args: args{
				ctx:     context.Background(),
				userId:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				trackId: uuid.MustParse("5b60e78f-b465-4cd6-b5d3-15b188f47a6a"),
			},
			want: fmt.Errorf("can't dislike track in repository"),
			setup: func(a args, f fields) {
				f.repo.EXPECT().DislikeTrack(a.ctx, a.userId, a.trackId).Return(fmt.Errorf("can't dislike track in repository"))
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
			u := usecase.NewUserInteractor(f.repo)

			tt.setup(tt.args, f)

			err := u.DislikeTrack(tt.args.ctx, tt.args.userId, tt.args.trackId)
			if (err != nil) != tt.wantErr {
				t.Errorf("userInteractor.DislikeTrack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func MustParseTime(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

func Test_userInteractor_ShowLikedTracks(t *testing.T) {
	type fields struct {
		repo *repository.MockUserRepository
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		setup   func(a args, f fields)
		want    []*entity.MusicDB
		wantErr bool
	}{
		{
			name: "success: ShowLikedTracks userInteractor",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			setup: func(a args, f fields) {
				f.repo.EXPECT().ShowLikedTracks(a.ctx, a.id).Return([]*entity.MusicDB{
					{
						Id:       uuid.MustParse("5b60e78f-b465-4cd6-b5d3-15b188f47a6a"),
						Name:     "Song1",
						Release:  MustParseTime("2006-01-02", "2022-01-01"),
						FileName: "Song1",
						Size:     1000,
						Duration: "00:01:00",
					},
					{
						Id:       uuid.MustParse("5b60e89f-b465-4cd6-b5d3-15b188f47a6a"),
						Name:     "Song2",
						Release:  MustParseTime("2006-01-02", "2022-01-01"),
						FileName: "Song2",
						Size:     1000,
						Duration: "00:01:00",
					},
				}, nil)
			},
			want: []*entity.MusicDB{
				{
					Id:       uuid.MustParse("5b60e78f-b465-4cd6-b5d3-15b188f47a6a"),
					Name:     "Song1",
					Release:  MustParseTime("2006-01-02", "2022-01-01"),
					FileName: "Song1",
					Size:     1000,
					Duration: "00:01:00",
				},
				{
					Id:       uuid.MustParse("5b60e89f-b465-4cd6-b5d3-15b188f47a6a"),
					Name:     "Song2",
					Release:  MustParseTime("2006-01-02", "2022-01-01"),
					FileName: "Song2",
					Size:     1000,
					Duration: "00:01:00",
				},
			},
			wantErr: false,
		},
		{
			name: "error: ShowLikedTracks userInteractor",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			setup: func(a args, f fields) {
				f.repo.EXPECT().ShowLikedTracks(a.ctx, a.id).Return(nil, fmt.Errorf("can't show liked tracks in repository"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				repo: repository.NewMockUserRepository(ctrl),
			}
			u := usecase.NewUserInteractor(f.repo)

			tt.setup(tt.args, f)

			got, err := u.ShowLikedTracks(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("userInteractor.ShowLikedTracks() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userInteractor.ShowLikedTracks() = %v, want %v", got, tt.want)
			}
		})
	}
}
