package repository

import (
	"context"
	"fmt"
	"mime/multipart"
	"music-backend-test/internal/db"
	"music-backend-test/internal/entity"
	"music-backend-test/internal/repository"
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
		utils  *utils.MockMusicUtils
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
			defer ctrl.Finish()
			f := fields{
				source: db.NewMockMusicSource(ctrl),
				utils:  utils.NewMockMusicUtils(ctrl),
			}
			os := utils.NewMockOS()
			musicRepository := repository.NewMusicRepository(f.source, f.utils, os)

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
		utils  *utils.MockMusicUtils
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
				f.source.EXPECT().Get(a.ctx, a.musicId).Return(&entity.MusicDB{
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
				f.source.EXPECT().Get(a.ctx, a.musicId).Return(nil, fmt.Errorf("Error in source.GetAll()"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				source: db.NewMockMusicSource(ctrl),
				utils:  utils.NewMockMusicUtils(ctrl),
			}
			os := utils.NewMockOS()
			musicRepository := repository.NewMusicRepository(f.source, f.utils, os)

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
		utils  *utils.MockMusicUtils
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
				f.source.EXPECT().GetAndSortByPopular(a.ctx).Return([]*entity.MusicDB{
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
			defer ctrl.Finish()
			f := fields{
				source: db.NewMockMusicSource(ctrl),
				utils:  utils.NewMockMusicUtils(ctrl),
			}
			os := utils.NewMockOS()
			musicRepository := repository.NewMusicRepository(f.source, f.utils, os)

			tt.setup(tt.args, f)

			got, err := musicRepository.GetAndSortByPopular(tt.args.ctx)
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
		utils  *utils.MockMusicUtils
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
				f.source.EXPECT().GetAllSortByTime(a.ctx).Return([]*entity.MusicDB{
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
			defer ctrl.Finish()
			f := fields{
				source: db.NewMockMusicSource(ctrl),
				utils:  utils.NewMockMusicUtils(ctrl),
			}
			os := utils.NewMockOS()
			musicRepository := repository.NewMusicRepository(f.source, f.utils, os)

			tt.setup(tt.args, f)

			got, err := musicRepository.GetAllSortByTime(tt.args.ctx)
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
		utils  *utils.MockMusicUtils
	}
	type args struct {
		ctx        context.Context
		musicParse *entity.MusicParse
	}
	ctx := context.Background()

	tests := []struct {
		name                      string
		args                      args
		setupGetSupportedFileType func(a args, f fields) utils.FileType
		setupGetAudioDuration     func(fileType utils.FileType, filePath string, os utils.FileSystem, f fields) string
		setupCreate               func(ctx context.Context, musicCreate *entity.MusicDB, f fields)
		wantErr                   bool
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
						Size:     900,
					},
				},
			},
			setupGetSupportedFileType: func(a args, f fields) utils.FileType {
				f.utils.EXPECT().GetSupportedFileType(a.musicParse.FileHeader.Filename).Return(utils.FileType(utils.MP3), nil)
				return utils.FileType(utils.MP3)
			},
			setupGetAudioDuration: func(fileType utils.FileType, filePath string, os utils.FileSystem, f fields) string {
				f.utils.EXPECT().GetAudioDuration(fileType, filePath, os).Return("3:15", nil)
				return "3:15"
			},
			setupCreate: func(ctx context.Context, musicCreate *entity.MusicDB, f fields) {
				f.source.EXPECT().Create(ctx, musicCreate).Return(nil)
			},
			wantErr: false,
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
						Filename: "Test.OGG",
						Size:     900,
					},
				},
			},
			setupGetSupportedFileType: func(a args, f fields) utils.FileType {
				f.utils.EXPECT().GetSupportedFileType(a.musicParse.FileHeader.Filename).Return(utils.FileType(utils.Invalid), fmt.Errorf("Incorrect file type"))
				return utils.FileType(utils.Invalid)
			},
			wantErr: true,
		},
		{
			name: "Get error in GetAudioDuration",
			args: args{
				ctx: ctx,
				musicParse: &entity.MusicParse{
					Name:    "Song2",
					Release: time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
					File:    os.NewFile(uintptr(syscall.Stdout), "Test.MP3"),
					FileHeader: &multipart.FileHeader{
						Filename: "Test.MP3",
						Size:     900,
					},
				},
			},
			setupGetSupportedFileType: func(a args, f fields) utils.FileType {
				f.utils.EXPECT().GetSupportedFileType(a.musicParse.FileHeader.Filename).Return(utils.FileType(utils.MP3), nil)
				return utils.FileType(utils.MP3)
			},
			setupGetAudioDuration: func(fileType utils.FileType, filePath string, os utils.FileSystem, f fields) string {
				f.utils.EXPECT().GetAudioDuration(fileType, filePath, os).Return("", fmt.Errorf("Error in GetAudioDuration"))
				return ""
			},
			wantErr: true,
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
						Size:     900,
					},
				},
			},
			setupGetSupportedFileType: func(a args, f fields) utils.FileType {
				f.utils.EXPECT().GetSupportedFileType(a.musicParse.FileHeader.Filename).Return(utils.FileType(utils.MP3), nil)
				return utils.FileType(utils.MP3)
			},
			setupGetAudioDuration: func(fileType utils.FileType, filePath string, os utils.FileSystem, f fields) string {
				f.utils.EXPECT().GetAudioDuration(fileType, filePath, os).Return("3:15", nil)
				return "3:15"
			},
			setupCreate: func(ctx context.Context, musicCreate *entity.MusicDB, f fields) {
				f.source.EXPECT().Create(ctx, musicCreate).Return(fmt.Errorf("Error in source.GetAll()"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				source: db.NewMockMusicSource(ctrl),
				utils:  utils.NewMockMusicUtils(ctrl),
			}
			os := utils.NewMockOS()
			musicRepository := repository.NewMusicRepository(f.source, f.utils, os)

			fileType := tt.setupGetSupportedFileType(tt.args, f)
			var duration string
			if tt.setupGetAudioDuration != nil {
				duration = tt.setupGetAudioDuration(fileType, "./internal/storage/music_storage/"+tt.args.musicParse.FileHeader.Filename, musicRepository.FileSystem, f)
			}

			musicDB := &entity.MusicDB{
				Name:     tt.args.musicParse.Name,
				Release:  tt.args.musicParse.Release,
				FileName: tt.args.musicParse.FileHeader.Filename,
				Size:     uint64(tt.args.musicParse.FileHeader.Size),
				Duration: duration,
			}
			if tt.setupCreate != nil {
				tt.setupCreate(tt.args.ctx, musicDB, f)
			}

			err := musicRepository.Create(tt.args.ctx, tt.args.musicParse)
			if tt.wantErr == true {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_Update(t *testing.T) {
	type fields struct {
		source *db.MockMusicSource
		utils  *utils.MockMusicUtils
	}
	type args struct {
		ctx        context.Context
		id         uuid.UUID
		musicParse *entity.MusicParse
	}
	ctx := context.Background()

	tests := []struct {
		name                      string
		args                      args
		setupGetSupportedFileType func(a args, f fields) utils.FileType
		setupGetAudioDuration     func(fileType utils.FileType, filePath string, os utils.FileSystem, f fields) string
		setupGet                  func(a args, f fields)
		setupUpdate               func(ctx context.Context, musicCreate *entity.MusicDB, f fields)
		wantErr                   bool
	}{
		{
			name: "Update music",
			args: args{
				ctx: ctx,
				id:  uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
				musicParse: &entity.MusicParse{
					Name:    "Song2",
					Release: time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
					File:    os.NewFile(uintptr(syscall.Stdout), "Test.MP3"),
					FileHeader: &multipart.FileHeader{
						Filename: "Test.MP3",
						Size:     900,
					},
				},
			},
			setupGetSupportedFileType: func(a args, f fields) utils.FileType {
				f.utils.EXPECT().GetSupportedFileType(a.musicParse.FileHeader.Filename).Return(utils.FileType(utils.MP3), nil)
				return utils.FileType(utils.MP3)
			},
			setupGet: func(a args, f fields) {
				f.source.EXPECT().Get(a.ctx, a.id).Return(&entity.MusicDB{
					Id:       uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
					Name:     "Song2",
					Release:  time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
					FileName: "Test.MP3",
					Size:     uint64(500),
					Duration: "2:47",
				}, nil)
			},
			setupGetAudioDuration: func(fileType utils.FileType, filePath string, os utils.FileSystem, f fields) string {
				f.utils.EXPECT().GetAudioDuration(fileType, filePath, os).Return("3:15", nil)
				return "3:15"
			},
			setupUpdate: func(ctx context.Context, musicUpdate *entity.MusicDB, f fields) {
				f.source.EXPECT().Update(ctx, musicUpdate).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Update music witout file",
			args: args{
				ctx: ctx,
				id:  uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
				musicParse: &entity.MusicParse{
					Name:    "Song2",
					Release: time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
				},
			},
			setupUpdate: func(ctx context.Context, musicUpdate *entity.MusicDB, f fields) {
				f.source.EXPECT().Update(ctx, musicUpdate).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Incorrect file type",
			args: args{
				ctx: ctx,
				id:  uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
				musicParse: &entity.MusicParse{
					Name:    "Song2",
					Release: time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
					File:    os.NewFile(uintptr(syscall.Stdout), "Test.MP3"),
					FileHeader: &multipart.FileHeader{
						Filename: "Test.OGG",
						Size:     900,
					},
				},
			},
			setupGetSupportedFileType: func(a args, f fields) utils.FileType {
				f.utils.EXPECT().GetSupportedFileType(a.musicParse.FileHeader.Filename).Return(utils.FileType(utils.Invalid), fmt.Errorf("Incorrect file type"))
				return utils.FileType(utils.Invalid)
			},
			wantErr: true,
		},
		{
			name: "Get error in GetAudioDuration",
			args: args{
				ctx: ctx,
				id:  uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
				musicParse: &entity.MusicParse{
					Name:    "Song2",
					Release: time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
					File:    os.NewFile(uintptr(syscall.Stdout), "Test.MP3"),
					FileHeader: &multipart.FileHeader{
						Filename: "Test.MP3",
						Size:     900,
					},
				},
			},
			setupGetSupportedFileType: func(a args, f fields) utils.FileType {
				f.utils.EXPECT().GetSupportedFileType(a.musicParse.FileHeader.Filename).Return(utils.FileType(utils.MP3), nil)
				return utils.FileType(utils.MP3)
			},
			setupGet: func(a args, f fields) {
				f.source.EXPECT().Get(a.ctx, a.id).Return(&entity.MusicDB{
					Id:       uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
					Name:     "Song2",
					Release:  time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
					FileName: "Test.MP3",
					Size:     uint64(500),
					Duration: "2:47",
				}, nil)
			},
			setupGetAudioDuration: func(fileType utils.FileType, filePath string, os utils.FileSystem, f fields) string {
				f.utils.EXPECT().GetAudioDuration(fileType, filePath, os).Return("", fmt.Errorf("Error in GetAudioDuration"))
				return ""
			},
			wantErr: true,
		},
		{
			name: "Get error from source",
			args: args{
				ctx: ctx,
				id:  uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
				musicParse: &entity.MusicParse{
					Name:    "Song2",
					Release: time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
					File:    os.NewFile(uintptr(syscall.Stdout), "Test.MP3"),
					FileHeader: &multipart.FileHeader{
						Filename: "Test.MP3",
						Size:     900,
					},
				},
			},
			setupGetSupportedFileType: func(a args, f fields) utils.FileType {
				f.utils.EXPECT().GetSupportedFileType(a.musicParse.FileHeader.Filename).Return(utils.FileType(utils.MP3), nil)
				return utils.FileType(utils.MP3)
			},
			setupGet: func(a args, f fields) {
				f.source.EXPECT().Get(a.ctx, a.id).Return(&entity.MusicDB{
					Id:       uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
					Name:     "Song2",
					Release:  time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
					FileName: "Test.MP3",
					Size:     uint64(500),
					Duration: "2:47",
				}, nil)
			},
			setupGetAudioDuration: func(fileType utils.FileType, filePath string, os utils.FileSystem, f fields) string {
				f.utils.EXPECT().GetAudioDuration(fileType, filePath, os).Return("3:15", nil)
				return "3:15"
			},
			setupUpdate: func(ctx context.Context, musicCreate *entity.MusicDB, f fields) {
				f.source.EXPECT().Update(ctx, musicCreate).Return(fmt.Errorf("Error in source.GetAll()"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				source: db.NewMockMusicSource(ctrl),
				utils:  utils.NewMockMusicUtils(ctrl),
			}
			os := utils.NewMockOS()
			musicRepository := repository.NewMusicRepository(f.source, f.utils, os)

			if tt.setupGet != nil {
				tt.setupGet(tt.args, f)
			}
			var fileType utils.FileType
			if tt.setupGetSupportedFileType != nil {
				fileType = tt.setupGetSupportedFileType(tt.args, f)
			}
			var duration string
			if tt.setupGetAudioDuration != nil {
				duration = tt.setupGetAudioDuration(fileType, "./internal/storage/music_storage/"+tt.args.musicParse.FileHeader.Filename, musicRepository.FileSystem, f)
			}

			filename := ""
			filesize := int64(0)
			if tt.args.musicParse.FileHeader != nil {
				filename = tt.args.musicParse.FileHeader.Filename
				filesize = tt.args.musicParse.FileHeader.Size
			}
			musicDB := &entity.MusicDB{
				Id:       uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
				Name:     tt.args.musicParse.Name,
				Release:  tt.args.musicParse.Release,
				FileName: filename,
				Size:     uint64(filesize),
				Duration: duration,
			}
			if tt.setupUpdate != nil {
				tt.setupUpdate(tt.args.ctx, musicDB, f)
			}

			err := musicRepository.Update(tt.args.ctx, tt.args.id, tt.args.musicParse)
			if tt.wantErr == true {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_Delete(t *testing.T) {
	type fields struct {
		source *db.MockMusicSource
		utils  *utils.MockMusicUtils
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
			setupGet: func(a args, f fields) {
				f.source.EXPECT().Get(a.ctx, a.musicId).Return(&entity.MusicDB{
					Id:       uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
					Name:     "Song1",
					Release:  time.Date(2023, time.March, 24, 0, 0, 0, 0, time.UTC),
					FileName: "Song1.mp3",
					Size:     uint64(500),
					Duration: "2:47",
				}, nil)
			},
			setupDelete: func(a args, f fields) {
				f.source.EXPECT().Delete(a.ctx, a.musicId).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Not in base",
			args: args{
				ctx:     context.Background(),
				musicId: uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
			},
			setupGet: func(a args, f fields) {
				f.source.EXPECT().Get(a.ctx, a.musicId).Return(nil, fmt.Errorf("Not in base"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				source: db.NewMockMusicSource(ctrl),
				utils:  utils.NewMockMusicUtils(ctrl),
			}
			os := utils.NewMockOS()
			musicSource := repository.NewMusicRepository(f.source, f.utils, os)
			tt.setupGet(tt.args, f)
			if tt.setupDelete != nil {
				tt.setupDelete(tt.args, f)
			}

			err := musicSource.Delete(tt.args.ctx, tt.args.musicId)
			if tt.wantErr == true {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
