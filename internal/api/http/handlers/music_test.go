package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"music-backend-test/internal/api/http/presenter"
	"music-backend-test/internal/entity"
	"music-backend-test/internal/usecase"
	"net/http"
	"net/http/httptest"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_GetAll(t *testing.T) {
	type fields struct {
		usecase *usecase.MockMusicInteractor
	}
	ctx := context.Background()

	tests := []struct {
		name       string
		setup      func(ctx context.Context, f fields)
		wantStatus int
		wantBody   string
	}{
		{
			name: "GetAll",
			setup: func(ctx context.Context, f fields) {
				f.usecase.EXPECT().GetAll(ctx).Return([]*entity.MusicDB{
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
			wantStatus: http.StatusOK,
			wantBody:   `[{"id":"4a6e104d-9d7f-45ff-8de6-37993d709522","name":"Song1","size":"500 B","duration":"2:47"},{"id":"ff578289-cdca-406e-9a57-f8c773f0cd15","name":"Song2","size":"900 B","duration":"3:23"}]`,
		},
		{
			name: "Error in usecase GetAll",
			setup: func(ctx context.Context, f fields) {
				f.usecase.EXPECT().GetAll(ctx).Return(nil, fmt.Errorf("Error in usecase GetAll"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntr := gomock.NewController(t)
			defer cntr.Finish()
			f := fields{
				usecase: usecase.NewMockMusicInteractor(cntr),
			}
			tt.setup(ctx, f)

			musicHandler := NewMusicHandlers(f.usecase, presenter.NewPresenter())

			r := gin.New()
			r.POST("/getall", musicHandler.GetAll)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/getall", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}

func Test_Get(t *testing.T) { //Очень нужен код ревью
	type fields struct {
		usecase *usecase.MockMusicInteractor
	}
	ctx := context.Background()

	tests := []struct {
		name       string
		setup      func(ctx context.Context, musicId uuid.UUID, f fields)
		id         string
		filename   string
		wantStatus int
	}{
		{
			name: "GetAll",
			setup: func(ctx context.Context, musicId uuid.UUID, f fields) {
				f.usecase.EXPECT().Get(ctx, musicId).Return(&entity.MusicDB{
					Id:       uuid.MustParse("ff578289-cdca-406e-9a57-f8c773f0cd15"),
					Name:     "Song2",
					Release:  time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
					FileName: "Song2.mp3",
					Size:     uint64(900),
					Duration: "3:23",
				}, nil)
			},
			id:         "ff578289-cdca-406e-9a57-f8c773f0cd15",
			filename:   "Song2.mp3",
			wantStatus: 0,
		},
		{
			name:       "Parse id error",
			id:         "0",
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Error in usecase GetAll",
			setup: func(ctx context.Context, musicId uuid.UUID, f fields) {
				f.usecase.EXPECT().Get(ctx, musicId).Return(nil, fmt.Errorf("Error in usecase Get"))
			},
			id:         "ff578289-cdca-406e-9a57-f8c773f0cd15",
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntr := gomock.NewController(t)
			defer cntr.Finish()
			f := fields{
				usecase: usecase.NewMockMusicInteractor(cntr),
			}
			if tt.setup != nil {
				tt.setup(ctx, uuid.MustParse(tt.id), f)
			}

			musicHandler := NewMusicHandlers(f.usecase, presenter.NewPresenter())

			r := gin.New()
			r.POST("/get/:id", musicHandler.Get)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/get/"+tt.id, nil)

			r.ServeHTTP(w, req)

			if tt.wantStatus == 0 {
				assert.Contains(t, w.Header()["Content-Disposition"], `attachment; filename="`+tt.filename+`"`)
			} else {
				assert.Equal(t, tt.wantStatus, w.Code)
			}

		})
	}
}

func Test_GetAllSortByTime(t *testing.T) {
	type fields struct {
		usecase *usecase.MockMusicInteractor
	}
	ctx := context.Background()

	tests := []struct {
		name       string
		setup      func(ctx context.Context, f fields)
		wantStatus int
		wantBody   string
	}{
		{
			name: "GetAllSortByTime",
			setup: func(ctx context.Context, f fields) {
				f.usecase.EXPECT().GetAllSortByTime(ctx).Return([]*entity.MusicDB{
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
			wantStatus: http.StatusOK,
			wantBody:   `[{"id":"4a6e104d-9d7f-45ff-8de6-37993d709522","name":"Song1","size":"500 B","duration":"2:47"},{"id":"ff578289-cdca-406e-9a57-f8c773f0cd15","name":"Song2","size":"900 B","duration":"3:23"}]`,
		},
		{
			name: "Error in usecase GetAllSortByTime",
			setup: func(ctx context.Context, f fields) {
				f.usecase.EXPECT().GetAllSortByTime(ctx).Return(nil, fmt.Errorf("Error in usecase GetAllSortByTime"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntr := gomock.NewController(t)
			defer cntr.Finish()
			f := fields{
				usecase: usecase.NewMockMusicInteractor(cntr),
			}
			tt.setup(ctx, f)

			musicHandler := NewMusicHandlers(f.usecase, presenter.NewPresenter())

			r := gin.New()
			r.POST("/getSortByTime", musicHandler.GetAllSortByTime)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/getSortByTime", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}

func Test_GetAndSortByPopular(t *testing.T) {
	type fields struct {
		usecase *usecase.MockMusicInteractor
	}
	ctx := context.Background()

	tests := []struct {
		name       string
		setup      func(ctx context.Context, f fields)
		wantStatus int
		wantBody   string
	}{
		{
			name: "GetAndSortByPopular",
			setup: func(ctx context.Context, f fields) {
				f.usecase.EXPECT().GetAndSortByPopular(ctx).Return([]*entity.MusicDB{
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
			wantStatus: http.StatusOK,
			wantBody:   `[{"id":"4a6e104d-9d7f-45ff-8de6-37993d709522","name":"Song1","size":"500 B","duration":"2:47"},{"id":"ff578289-cdca-406e-9a57-f8c773f0cd15","name":"Song2","size":"900 B","duration":"3:23"}]`,
		},
		{
			name: "Error in usecase GetAndSortByPopular",
			setup: func(ctx context.Context, f fields) {
				f.usecase.EXPECT().GetAndSortByPopular(ctx).Return(nil, fmt.Errorf("Error in usecase GetAndSortByPopular"))
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   "[]", //КАК?!?!
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntr := gomock.NewController(t)
			defer cntr.Finish()
			f := fields{
				usecase: usecase.NewMockMusicInteractor(cntr),
			}
			tt.setup(ctx, f)

			musicHandler := NewMusicHandlers(f.usecase, presenter.NewPresenter())

			r := gin.New()
			r.POST("/getSortByPopular", musicHandler.GetAndSortByPopular)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/getSortByPopular", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}

func Test_Create(t *testing.T) {
	type fields struct {
		usecase *usecase.MockMusicInteractor
	}
	ctx := context.Background()
	tests := []struct {
		name       string
		setup      func(ctx context.Context, musicCreate entity.MusicParse, f fields)
		inputBody  func(w *multipart.Writer)
		wantStatus int
	}{
		{
			name: "Create",
			setup: func(ctx context.Context, musicCreate entity.MusicParse, f fields) {
				f.usecase.EXPECT().Create(ctx, musicCreate).Return(nil)
			},
			inputBody: func(w *multipart.Writer) {
				name, err := w.CreateFormField("name")
				assert.NoError(t, err)
				name.Write([]byte("Song2"))

				release, err := w.CreateFormField("release")
				assert.NoError(t, err)
				release.Write([]byte("2021-11-15"))

				file, err := w.CreateFormFile("file", "Test.mp3")
				assert.NoError(t, err)

				data, err := os.Open("./internal/storage/music_storage/Test.mp3") //Неверный путь
				assert.NoError(t, err)
				defer data.Close()
				_, err = io.Copy(file, data)
				assert.NoError(t, err)

				w.Close()
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Error in usecase create",
			setup: func(ctx context.Context, musicCreate entity.MusicParse, f fields) {
				f.usecase.EXPECT().Create(ctx, musicCreate).Return(fmt.Errorf("Error in usecase create"))
			},
			inputBody: func(w *multipart.Writer) {
				name, err := w.CreateFormField("name")
				assert.NoError(t, err)
				name.Write([]byte("Song2"))

				release, err := w.CreateFormField("release")
				assert.NoError(t, err)
				release.Write([]byte("2021-11-15"))

				file, err := w.CreateFormFile("file", "Test.mp3")
				assert.NoError(t, err)

				data, err := os.Open("./internal/storage/music_storage/Test.mp3") //Неверный путь
				assert.NoError(t, err)
				defer data.Close()
				_, err = io.Copy(file, data)
				assert.NoError(t, err)

				w.Close()
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				usecase: usecase.NewMockMusicInteractor(ctrl),
			}

			musicHandler := NewMusicHandlers(f.usecase, presenter.NewPresenter())
			tt.setup(ctx, entity.MusicParse{ //Несостыковка приходящих и статичных данных
				Name:    "Song2",
				Release: time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
				File:    os.NewFile(uintptr(syscall.Stdout), "Test.mp3"),
				FileHeader: &multipart.FileHeader{
					Filename: "Test.mp3",
					Size:     64,
				},
			}, f)

			r := gin.New()
			r.POST("/create", musicHandler.Create)

			buf := new(bytes.Buffer)
			writer := multipart.NewWriter(buf)
			tt.inputBody(writer)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/create", buf)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func Update(t *testing.T) {
	type fields struct {
		usecase *usecase.MockMusicInteractor
	}
	ctx := context.Background()
	tests := []struct {
		name       string
		setup      func(ctx context.Context, musicId uuid.UUID, music entity.MusicParse, f fields)
		id         string
		inputBody  func(w *multipart.Writer)
		wantStatus int
	}{
		{
			name: "Update",
			setup: func(ctx context.Context, musicId uuid.UUID, music entity.MusicParse, f fields) {
				f.usecase.EXPECT().Update(ctx, musicId, music).Return(nil)
			},
			id: "ff578289-cdca-406e-9a57-f8c773f0cd15",
			inputBody: func(w *multipart.Writer) {
				name, err := w.CreateFormField("name")
				assert.NoError(t, err)
				name.Write([]byte("Song2"))

				release, err := w.CreateFormField("release")
				assert.NoError(t, err)
				release.Write([]byte("2021-11-15"))

				file, err := w.CreateFormFile("file", "Test.mp3")
				assert.NoError(t, err)

				data, err := os.Open("./internal/storage/music_storage/Test.mp3") //Неверный путь
				assert.NoError(t, err)
				defer data.Close()
				_, err = io.Copy(file, data)
				assert.NoError(t, err)

				w.Close()
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "Error in id parse",
			id:         "0",
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Error in usecase Update",
			setup: func(ctx context.Context, musicId uuid.UUID, music entity.MusicParse, f fields) {
				f.usecase.EXPECT().Update(ctx, musicId, music).Return(fmt.Errorf("Error in usecase Update"))
			},
			id: "ff578289-cdca-406e-9a57-f8c773f0cd15",
			inputBody: func(w *multipart.Writer) {
				name, err := w.CreateFormField("name")
				assert.NoError(t, err)
				name.Write([]byte("Song2"))

				release, err := w.CreateFormField("release")
				assert.NoError(t, err)
				release.Write([]byte("2021-11-15"))

				file, err := w.CreateFormFile("file", "Test.mp3")
				assert.NoError(t, err)

				data, err := os.Open("./internal/storage/music_storage/Test.mp3") //Неверный путь
				assert.NoError(t, err)
				defer data.Close()
				_, err = io.Copy(file, data)
				assert.NoError(t, err)

				w.Close()
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				usecase: usecase.NewMockMusicInteractor(ctrl),
			}

			musicHandler := NewMusicHandlers(f.usecase, presenter.NewPresenter())
			if tt.setup != nil {
				tt.setup(ctx, uuid.MustParse(tt.id), entity.MusicParse{ //Несостыковка приходящих и статичных данных
					Name:    "Song2",
					Release: time.Date(2021, time.November, 15, 0, 0, 0, 0, time.UTC),
					File:    os.NewFile(uintptr(syscall.Stdout), "Test.mp3"),
					FileHeader: &multipart.FileHeader{
						Filename: "Test.mp3",
						Size:     64,
					}}, f)
			}

			r := gin.New()
			r.POST("/update/:id", musicHandler.Update)

			buf := new(bytes.Buffer)
			writer := multipart.NewWriter(buf)
			if tt.inputBody != nil {
				tt.inputBody(writer)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/update/"+tt.id, buf)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func Test_Delete(t *testing.T) {
	type fields struct {
		usecase *usecase.MockMusicInteractor
	}
	ctx := context.Background()

	tests := []struct {
		name       string
		setup      func(ctx context.Context, musicId uuid.UUID, f fields)
		id         string
		wantStatus int
	}{
		{
			name: "Delete",
			setup: func(ctx context.Context, musicId uuid.UUID, f fields) {
				f.usecase.EXPECT().Delete(ctx, musicId).Return(nil)
			},
			id:         "ff578289-cdca-406e-9a57-f8c773f0cd15",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Error in id parse",
			id:         "0",
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Error in usecase Delete",
			setup: func(ctx context.Context, musicId uuid.UUID, f fields) {
				f.usecase.EXPECT().Delete(ctx, musicId).Return(fmt.Errorf("Error in usecase Get"))
			},
			id:         "ff578289-cdca-406e-9a57-f8c773f0cd15",
			wantStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntr := gomock.NewController(t)
			defer cntr.Finish()
			f := fields{
				usecase: usecase.NewMockMusicInteractor(cntr),
			}
			if tt.setup != nil {
				tt.setup(ctx, uuid.MustParse(tt.id), f)
			}

			musicHandler := NewMusicHandlers(f.usecase, presenter.NewPresenter())

			r := gin.New()
			r.POST("/delete/:id", musicHandler.Delete)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/delete/"+tt.id, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
