package repository

import (
	"context"
	"fmt"
	"mime/multipart"
	"music-backend-test/internal/db"
	"music-backend-test/internal/entity"
	"music-backend-test/internal/utils"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_GetAll(t *testing.T) {
	type fields struct {
		source *db.MockMusicSource
	}
	type args struct {
		ctx context.Context
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		args    args
		setup   func(a args, f fields)
		want    []*entity.MusicDB
		wantErr bool
	}{
		{
			name: "Get all music",
			args: args{
				ctx: ctx,
			},
			setup: func(a args, f fields) {
				f.source.EXPECT().GetAll(a.ctx).Return([]*entity.MusicDB{
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
					}}, nil)
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
			name: "Get error from source",
			args: args{
				ctx: ctx,
			},
			setup: func(a args, f fields) {
				f.source.EXPECT().GetAll(a.ctx).Return(nil, fmt.Errorf("Error in source.GetAll()"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				source: db.NewMockMusicSource(ctrl),
			}

			musicRepository := NewMusicRepository(f.source)

			tt.setup(tt.args, f)

			got, err := musicRepository.GetAll(tt.args.ctx)
			if tt.wantErr == true {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, got)
				}
			}
		})

	}
}

func Test_Get(t *testing.T) {
	type fields struct {
		source *db.MockMusicSource
	}
	type args struct {
		ctx     context.Context
		musicId uuid.UUID
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		args    args
		setup   func(a args, f fields)
		want    *entity.MusicDB
		wantErr bool
	}{
		{
			name: "Get track by id",
			args: args{
				ctx:     context.Background(),
				musicId: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			setup: func(a args, f fields) {
				f.source.EXPECT().GetAll(a.ctx).Return(&entity.MusicDB{
					Id:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
					Name:     "Song1",
					Release:  time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
					FileName: "Song1.mp3",
					Size:     uint64(500),
					Duration: "2:47",
				}, nil)
			},
			want: &entity.MusicDB{
				Id:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				Name:     "Song1",
				Release:  time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
				FileName: "Song1.mp3",
				Size:     uint64(500),
				Duration: "2:47",
			},
			wantErr: false,
		},
		{
			name: "Get error from source",
			args: args{
				ctx: ctx,
			},
			setup: func(a args, f fields) {
				f.source.EXPECT().GetAll(a.ctx).Return(nil, fmt.Errorf("Error in source.GetAll()"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				source: db.NewMockMusicSource(ctrl),
			}

			musicRepository := NewMusicRepository(f.source)

			tt.setup(tt.args, f)

			got, err := musicRepository.Get(tt.args.ctx, tt.args.musicId)
			if tt.wantErr == true {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, got)
				}
			}
		})

	}
}

func Test_GetAndSortByPopular(t *testing.T) {
	type fields struct {
		source *db.MockMusicSource
	}
	type args struct {
		ctx context.Context
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		args    args
		setup   func(a args, f fields)
		want    []*entity.MusicDB
		wantErr bool
	}{
		{
			name: "Get all music",
			args: args{
				ctx: ctx,
			},
			setup: func(a args, f fields) {
				f.source.EXPECT().GetAll(a.ctx).Return([]*entity.MusicDB{
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
					}}, nil)
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
			name: "Get error from source",
			args: args{
				ctx: ctx,
			},
			setup: func(a args, f fields) {
				f.source.EXPECT().GetAndSortByPopular(a.ctx).Return(nil, fmt.Errorf("Error in source.GetAll()"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				source: db.NewMockMusicSource(ctrl),
			}

			musicRepository := NewMusicRepository(f.source)

			tt.setup(tt.args, f)

			got, err := musicRepository.GetAll(tt.args.ctx)
			if tt.wantErr == true {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, got)
				}
			}
		})

	}
}

func Test_GetAllSortByTime(t *testing.T) {
	type fields struct {
		source *db.MockMusicSource
	}
	type args struct {
		ctx context.Context
	}
	ctx := context.Background()

	tests := []struct {
		name    string
		args    args
		setup   func(a args, f fields)
		want    []*entity.MusicDB
		wantErr bool
	}{
		{
			name: "Get all music",
			args: args{
				ctx: ctx,
			},
			setup: func(a args, f fields) {
				f.source.EXPECT().GetAll(a.ctx).Return([]*entity.MusicDB{
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
					}}, nil)
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
			name: "Get error from source",
			args: args{
				ctx: ctx,
			},
			setup: func(a args, f fields) {
				f.source.EXPECT().GetAllSortByTime(a.ctx).Return(nil, fmt.Errorf("Error in source.GetAll()"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				source: db.NewMockMusicSource(ctrl),
			}

			musicRepository := NewMusicRepository(f.source)

			tt.setup(tt.args, f)

			got, err := musicRepository.GetAll(tt.args.ctx)
			if tt.wantErr == true {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, got)
				}
			}
		})

	}
}

func Test_Create(t *testing.T) {
	type fields struct {
		source *db.MockMusicSource
	}
	type args struct {
		ctx        context.Context
		musicParse *entity.MusicParse
		fileType   utils.FileType
	}
	ctx := context.Background()

	tests := []struct {
		name                  string
		args                  args
		setup                 func(ctx context.Context, musicCreate *entity.MusicDB, f fields)
		setupGetAudioDuration func(fileType utils.FileType, filePath string, f fields)
		wantErr               bool
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
				fileType: "MP3",
			},
			setup: func(ctx context.Context, musicCreate *entity.MusicDB, f fields) {
				f.source.EXPECT().Create(ctx, musicCreate).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Get error from source",
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
			setup: func(ctx context.Context, musicCreate *entity.MusicDB, f fields) {
				f.source.EXPECT().Create(ctx, musicCreate).Return(fmt.Errorf("Error in source.GetAll()"))
			},
			wantErr: true,
		},
		{
			name: "Incorrect file type",
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
				fileType: "OGG",
			},
			setup: func(ctx context.Context, musicCreate *entity.MusicDB, f fields) {
				f.source.EXPECT().Create(ctx, musicCreate).Return(nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			f := fields{
				source: db.NewMockMusicSource(ctrl),
			}

			musicRepository := NewMusicRepository(f.source)

			musicDB := &entity.MusicDB{
				Id:       uuid.Nil,
				Name:     tt.args.musicParse.Name,
				Release:  tt.args.musicParse.Release,
				FileName: "Song2.mp3",
				Size:     uint64(900),
				Duration: "3:23",
			}
			tt.setup(tt.args.ctx, musicDB, f)

			err := musicRepository.Create(tt.args.ctx, tt.args.musicParse, tt.args.fileType)
			if tt.wantErr == true {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})

	}
}

// func Test_Update(t *testing.T)

func Test_Delete(t *testing.T) {
	type fields struct {
		source *db.MockMusicSource
	}

	type args struct {
		ctx     context.Context
		musicId uuid.UUID
	}

	tests := []struct {
		name        string
		args        args
		setupDelete func(a args, f fields)
		setupGet    func(a args, f fields)
		wantErr     bool
	}{
		{
			name: "Delete music",
			args: args{
				ctx:     context.Background(),
				musicId: uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
			},
			setupDelete: func(a args, f fields) {
				f.source.EXPECT().Delete(a.ctx, a.musicId).Return(nil)
			},
			setupGet: func(a args, f fields) {
				f.source.EXPECT().Get(a.ctx, a.musicId).Return(&entity.MusicDB{
					Id:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
					Name:     "Song1",
					Release:  time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
					FileName: "Song1.mp3",
					Size:     uint64(500),
					Duration: "2:47",
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "Bad request to database at music.Delete",
			args: args{
				ctx:     context.Background(),
				musicId: uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
			},
			setupDelete: func(a args, f fields) {
				f.source.EXPECT().Delete(a.ctx, a.musicId).Return(nil)
			},
			setupGet: func(a args, f fields) {
				f.source.EXPECT().Get(a.ctx, a.musicId).Return(&entity.MusicDB{
					Id:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
					Name:     "Song1",
					Release:  time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
					FileName: "Song1.mp3",
					Size:     uint64(500),
					Duration: "2:47",
				}, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntr := gomock.NewController(t)
			f := fields{
				source: db.NewMockMusicSource(cntr),
			}

			musicSource := NewMusicRepository(f.source)
			tt.setupGet(tt.args, f)
			tt.setupDelete(tt.args, f)

			err := musicSource.Delete(tt.args.ctx, tt.args.musicId)
			if tt.wantErr == true {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
