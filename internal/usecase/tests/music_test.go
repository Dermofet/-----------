package usecase

import (
	"context"
	"fmt"
	"mime/multipart"
	"music-backend-test/internal/entity"
	"music-backend-test/internal/repository"
	"music-backend-test/internal/usecase"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_GetAll(t *testing.T) {
	type field struct {
		repository *repository.MockMusicRepository
	}

	type args struct {
		ctx context.Context
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		args    args
		setup   func(a args, f field)
		want    []*entity.MusicDB
		wantErr bool
	}{
		{
			name: "GetAll",
			args: args{
				ctx: ctx,
			},
			setup: func(a args, f field) {
				f.repository.EXPECT().GetAll(a.ctx).Return([]*entity.MusicDB{
					{
						Id:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
						Name:     "Song1",
						Release:  time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
						FileName: "Song1.mp3",
						Size:     uint64(500),
						Duration: "2:47",
					},
					{
						Id:       uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
						Name:     "Song2",
						Release:  time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
						FileName: "Song2.mp3",
						Size:     uint64(900),
						Duration: "3:23",
					},
				}, nil)
			},
			want: []*entity.MusicDB{
				{
					Id:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
					Name:     "Song1",
					Release:  time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
					FileName: "Song1.mp3",
					Size:     uint64(500),
					Duration: "2:47",
				},
				{
					Id:       uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
					Name:     "Song2",
					Release:  time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
					FileName: "Song2.mp3",
					Size:     uint64(900),
					Duration: "3:23",
				},
			},
			wantErr: false,
		},
		{
			name: "Error in repository GetAll",
			args: args{
				ctx: ctx,
			},
			setup: func(a args, f field) {
				f.repository.EXPECT().GetAll(a.ctx).Return(nil, fmt.Errorf("Error in GetAll"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntr := gomock.NewController(t)
			f := field{
				repository: repository.NewMockMusicRepository(cntr),
			}
			musicUsecase := usecase.NewMusicInteractor(f.repository)
			tt.setup(tt.args, f)

			got, gotErr := musicUsecase.GetAll(tt.args.ctx)
			if tt.wantErr == true {
				assert.Error(t, gotErr)
			} else {
				if assert.NoError(t, gotErr) {
					assert.Equal(t, tt.want, got)
				}
			}
		})
	}
}

func Test_Get(t *testing.T) {
	type field struct {
		repository *repository.MockMusicRepository
	}

	type args struct {
		ctx     context.Context
		musicId uuid.UUID
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		args    args
		setup   func(a args, f field)
		want    *entity.MusicDB
		wantErr bool
	}{
		{
			name: "Get music",
			args: args{
				ctx:     ctx,
				musicId: uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
			},
			setup: func(a args, f field) {
				f.repository.EXPECT().Get(a.ctx, a.musicId).Return(&entity.MusicDB{
					Id:       uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
					Name:     "Song2",
					Release:  time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
					FileName: "Song2.mp3",
					Size:     uint64(900),
					Duration: "3:23",
				}, nil)
			},
			want: &entity.MusicDB{
				Id:       uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
				Name:     "Song2",
				Release:  time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
				FileName: "Song2.mp3",
				Size:     uint64(900),
				Duration: "3:23",
			},
			wantErr: false,
		},
		{
			name: "Error in repository Get",
			args: args{
				ctx:     ctx,
				musicId: uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
			},
			setup: func(a args, f field) {
				f.repository.EXPECT().Get(a.ctx, a.musicId).Return(nil, fmt.Errorf("Error in repository Get"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntr := gomock.NewController(t)
			f := field{
				repository: repository.NewMockMusicRepository(cntr),
			}
			musicUsecase := usecase.NewMusicInteractor(f.repository)
			tt.setup(tt.args, f)

			got, gotErr := musicUsecase.Get(tt.args.ctx, tt.args.musicId)
			if tt.wantErr == true {
				assert.Error(t, gotErr)
			} else {
				if assert.NoError(t, gotErr) {
					assert.Equal(t, tt.want, got)
				}
			}
		})
	}
}
func Test_GetAndSortByPopular(t *testing.T) {
	type field struct {
		repository *repository.MockMusicRepository
	}

	type args struct {
		ctx context.Context
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		args    args
		setup   func(a args, f field)
		want    []*entity.MusicDB
		wantErr bool
	}{
		{
			name: "GetAndSortByPopular",
			args: args{
				ctx: ctx,
			},
			setup: func(a args, f field) {
				f.repository.EXPECT().GetAndSortByPopular(a.ctx).Return([]*entity.MusicDB{
					{
						Id:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
						Name:     "Song1",
						Release:  time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
						FileName: "Song1.mp3",
						Size:     uint64(500),
						Duration: "2:47",
					},
					{
						Id:       uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
						Name:     "Song2",
						Release:  time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
						FileName: "Song2.mp3",
						Size:     uint64(900),
						Duration: "3:23",
					},
				}, nil)
			},
			want: []*entity.MusicDB{
				{
					Id:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
					Name:     "Song1",
					Release:  time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
					FileName: "Song1.mp3",
					Size:     uint64(500),
					Duration: "2:47",
				},
				{
					Id:       uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
					Name:     "Song2",
					Release:  time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
					FileName: "Song2.mp3",
					Size:     uint64(900),
					Duration: "3:23",
				},
			},
			wantErr: false,
		},
		{
			name: "Error in repository GetAndSortByPopular",
			args: args{
				ctx: ctx,
			},
			setup: func(a args, f field) {
				f.repository.EXPECT().GetAndSortByPopular(a.ctx).Return(nil, fmt.Errorf("Error in GetAndSortByPopular"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntr := gomock.NewController(t)
			f := field{
				repository: repository.NewMockMusicRepository(cntr),
			}
			musicUsecase := usecase.NewMusicInteractor(f.repository)
			tt.setup(tt.args, f)

			got, gotErr := musicUsecase.GetAndSortByPopular(tt.args.ctx)
			if tt.wantErr == true {
				assert.Error(t, gotErr)
			} else {
				if assert.NoError(t, gotErr) {
					assert.Equal(t, tt.want, got)
				}
			}
		})
	}
}

func Test_GetAllSortByTime(t *testing.T) {
	type field struct {
		repository *repository.MockMusicRepository
	}

	type args struct {
		ctx context.Context
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		args    args
		setup   func(a args, f field)
		want    []*entity.MusicDB
		wantErr bool
	}{
		{
			name: "GetAllSortByTime",
			args: args{
				ctx: ctx,
			},
			setup: func(a args, f field) {
				f.repository.EXPECT().GetAllSortByTime(a.ctx).Return([]*entity.MusicDB{
					{
						Id:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
						Name:     "Song1",
						Release:  time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
						FileName: "Song1.mp3",
						Size:     uint64(500),
						Duration: "2:47",
					},
					{
						Id:       uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
						Name:     "Song2",
						Release:  time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
						FileName: "Song2.mp3",
						Size:     uint64(900),
						Duration: "3:23",
					},
				}, nil)
			},
			want: []*entity.MusicDB{
				{
					Id:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
					Name:     "Song1",
					Release:  time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
					FileName: "Song1.mp3",
					Size:     uint64(500),
					Duration: "2:47",
				},
				{
					Id:       uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
					Name:     "Song2",
					Release:  time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
					FileName: "Song2.mp3",
					Size:     uint64(900),
					Duration: "3:23",
				},
			},
			wantErr: false,
		},
		{
			name: "Error in repository GetAllSortByTime",
			args: args{
				ctx: ctx,
			},
			setup: func(a args, f field) {
				f.repository.EXPECT().GetAllSortByTime(a.ctx).Return(nil, fmt.Errorf("Error in GetAllSortByTime"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntr := gomock.NewController(t)
			f := field{
				repository: repository.NewMockMusicRepository(cntr),
			}
			musicUsecase := usecase.NewMusicInteractor(f.repository)
			tt.setup(tt.args, f)

			got, gotErr := musicUsecase.GetAllSortByTime(tt.args.ctx)
			if tt.wantErr == true {
				assert.Error(t, gotErr)
			} else {
				if assert.NoError(t, gotErr) {
					assert.Equal(t, tt.want, got)
				}
			}
		})
	}
}

func Test_Create(t *testing.T) {
	type field struct {
		repository *repository.MockMusicRepository
	}

	type args struct {
		ctx        context.Context
		musicParse *entity.MusicParse
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		args    args
		setup   func(a args, f field)
		wantErr bool
	}{
		{
			name: "Create music",
			args: args{
				ctx: ctx,
				musicParse: &entity.MusicParse{
					Name:    "Song2",
					Release: time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
					File:    os.NewFile(uintptr(syscall.Stdout), "Test.MP3"),
					FileHeader: &multipart.FileHeader{
						Filename: "Test.MP3",
						Size:     64,
					},
				},
			},
			setup: func(a args, f field) {
				f.repository.EXPECT().Create(a.ctx, a.musicParse).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Error in repository Create",
			args: args{
				ctx: ctx,
				musicParse: &entity.MusicParse{
					Name:    "Song2",
					Release: time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
					File:    os.NewFile(uintptr(syscall.Stdout), "Test.MP3"),
					FileHeader: &multipart.FileHeader{
						Filename: "Test.MP3",
						Size:     64,
					},
				},
			},
			setup: func(a args, f field) {
				f.repository.EXPECT().Create(a.ctx, a.musicParse).Return(fmt.Errorf("Error in repository Create"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntr := gomock.NewController(t)
			f := field{
				repository: repository.NewMockMusicRepository(cntr),
			}
			musicUsecase := usecase.NewMusicInteractor(f.repository)
			tt.setup(tt.args, f)

			gotErr := musicUsecase.Create(tt.args.ctx, tt.args.musicParse)
			if tt.wantErr == true {
				assert.Error(t, gotErr)
			}
		})
	}
}

func Test_Update(t *testing.T) {
	type field struct {
		repository *repository.MockMusicRepository
	}

	type args struct {
		ctx        context.Context
		musicId    uuid.UUID
		parseMusic *entity.MusicParse
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		args    args
		setup   func(a args, f field)
		wantErr bool
	}{
		{
			name: "Update Music",
			args: args{
				ctx:     ctx,
				musicId: uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
				parseMusic: &entity.MusicParse{
					Name:    "Song2",
					Release: time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
					File:    os.NewFile(uintptr(syscall.Stdout), "Test.MP3"),
					FileHeader: &multipart.FileHeader{
						Filename: "Test.MP3",
						Size:     64,
					},
				},
			},
			setup: func(a args, f field) {
				f.repository.EXPECT().Update(a.ctx, a.musicId, a.parseMusic).Return(nil)
			},
		},
		{
			name: "Error in repositoty Update",
			args: args{
				ctx:     ctx,
				musicId: uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
				parseMusic: &entity.MusicParse{
					Name:    "Song2",
					Release: time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
					File:    os.NewFile(uintptr(syscall.Stdout), "Test.MP3"),
					FileHeader: &multipart.FileHeader{
						Filename: "Test.MP3",
						Size:     64,
					},
				},
			},
			setup: func(a args, f field) {
				f.repository.EXPECT().Update(a.ctx, a.musicId, a.parseMusic).Return(fmt.Errorf("Error in repositoty Update"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntr := gomock.NewController(t)
			f := field{
				repository: repository.NewMockMusicRepository(cntr),
			}
			musicUsecase := usecase.NewMusicInteractor(f.repository)
			tt.setup(tt.args, f)

			gotErr := musicUsecase.Update(tt.args.ctx, tt.args.musicId, tt.args.parseMusic)
			if tt.wantErr == true {
				assert.Error(t, gotErr)
			}
		})
	}
}

func Test_Delete(t *testing.T) {
	type field struct {
		repository *repository.MockMusicRepository
	}

	type args struct {
		ctx     context.Context
		musicId uuid.UUID
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		args    args
		setup   func(a args, f field)
		wantErr bool
	}{
		{
			name: "Delete music",
			args: args{
				ctx:     ctx,
				musicId: uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
			},
			setup: func(a args, f field) {
				f.repository.EXPECT().Delete(a.ctx, a.musicId).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Error in repository Delete",
			args: args{
				ctx:     ctx,
				musicId: uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
			},
			setup: func(a args, f field) {
				f.repository.EXPECT().Delete(a.ctx, a.musicId).Return(fmt.Errorf("Error in repository Delete"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntr := gomock.NewController(t)
			f := field{
				repository: repository.NewMockMusicRepository(cntr),
			}
			musicUsecase := usecase.NewMusicInteractor(f.repository)
			tt.setup(tt.args, f)

			gotErr := musicUsecase.Delete(tt.args.ctx, tt.args.musicId)
			if tt.wantErr == true {
				assert.Error(t, gotErr)
			}
		})
	}
}
