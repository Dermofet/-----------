package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"music-backend-test/internal/api/http/handlers"
	"music-backend-test/internal/api/http/presenter"
	"music-backend-test/internal/api/http/view"
	"music-backend-test/internal/entity"
	"music-backend-test/internal/usecase"
	"net/http"
	"net/http/httptest"
	reflect "reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func MustMarshal(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func Test_userHandlers_GetMeHandler(t *testing.T) {
	type fields struct {
		interactor *usecase.MockUserInteractor
		presenter  *presenter.MockPresenter
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	type сase struct {
		name           string
		args           args
		setup          func(a args, f fields)
		expectedStatus int
		expectedBody   *view.UserView
	}
	cases := []сase{
		{
			name: "GetMeHandler: 200",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			setup: func(a args, f fields) {
				userDB := &entity.UserDB{
					ID:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
					Username: "Paul",
					Password: "password",
					Role:     "user",
				}
				userView := &view.UserView{
					Id:       "4a6e104d-9d7f-45ff-8de6-37993d709522",
					Username: "Paul",
					Role:     "user",
				}

				f.interactor.EXPECT().GetById(a.ctx, a.id).Return(userDB, nil)
				f.presenter.EXPECT().ToUserView(userDB).Return(userView)
			},
			expectedStatus: 200,
			expectedBody: &view.UserView{
				Id:       "4a6e104d-9d7f-45ff-8de6-37993d709522",
				Username: "Paul",
				Role:     "user",
			},
		},
		{
			name: "GetMeHandler: 500",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			setup: func(a args, f fields) {
				f.interactor.EXPECT().GetById(a.ctx, a.id).Return(nil, fmt.Errorf("can't get user in interactor"))
			},
			expectedStatus: 500,
			expectedBody:   nil,
		},
		{
			name: "GetMeHandler: 404",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			setup: func(a args, f fields) {
				f.interactor.EXPECT().GetById(a.ctx, a.id).Return(nil, nil)
			},
			expectedStatus: 404,
			expectedBody:   nil,
		},
		{
			name: "GetMeHandler: 401",
			args: args{
				ctx: context.Background(),
				id:  uuid.Nil,
			},
			setup:          func(a args, f fields) {},
			expectedStatus: 401,
			expectedBody:   nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				interactor: usecase.NewMockUserInteractor(ctrl),
				presenter:  presenter.NewMockPresenter(ctrl),
			}

			h := handlers.NewUserHandlers(f.interactor, f.presenter)

			tc.setup(tc.args, f)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			if tc.args.id != uuid.Nil {
				c.Set("user-id", tc.args.id.String())
			}

			h.GetMeHandler(c)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			var userResponse *view.UserView

			body := w.Body.Bytes()
			if len(body) > 0 {
				fmt.Println(string(body))
				err := json.Unmarshal(w.Body.Bytes(), &userResponse)
				if err != nil {
					t.Errorf("can't unmarshal expected body: %v", err)
				}
			} else {
				userResponse = nil
			}

			if !reflect.DeepEqual(userResponse, tc.expectedBody) {
				t.Errorf("expected body %+v, got %+v", tc.expectedBody, userResponse)
			}
		})
	}
}

func Test_userHandlers_UpdateMeHandler(t *testing.T) {
	type fields struct {
		interactor *usecase.MockUserInteractor
		presenter  *presenter.MockPresenter
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	type testCase struct {
		name           string
		args           args
		setup          func(a args, f fields, c *gin.Context)
		expectedStatus int
		expectedBody   *view.UserView
	}
	cases := []testCase{
		{
			name: "UpdateMeHandler: 200",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			setup: func(a args, f fields, c *gin.Context) {
				userDB := &entity.UserDB{
					ID:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
					Username: "Paul",
					Password: "password",
					Role:     "admin",
				}
				userCreate := &entity.UserCreate{
					Username: "UpdatedPaul",
					Password: "updatedPassword",
				}
				userView := &view.UserView{
					Id:       "4a6e104d-9d7f-45ff-8de6-37993d709522",
					Username: "UpdatedPaul",
					Role:     "admin",
				}

				f.interactor.EXPECT().Update(a.ctx, a.id, userCreate).Return(userDB, nil)
				f.presenter.EXPECT().ToUserView(userDB).Return(userView)

				body, _ := json.Marshal(userCreate)
				c.Request = httptest.NewRequest("PUT", "/update", bytes.NewBuffer(body))
			},
			expectedStatus: 200,
			expectedBody: &view.UserView{
				Id:       "4a6e104d-9d7f-45ff-8de6-37993d709522",
				Username: "UpdatedPaul",
				Role:     "admin",
			},
		},
		{
			name: "UpdateMeHandler: 404",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			setup: func(a args, f fields, c *gin.Context) {
				userCreate := &entity.UserCreate{}

				f.interactor.EXPECT().Update(a.ctx, a.id, userCreate).Return(nil, nil)

				body, _ := json.Marshal(userCreate)
				c.Request = httptest.NewRequest("PUT", "/update", bytes.NewBuffer(body))
			},
			expectedStatus: 404,
			expectedBody:   nil,
		},
		{
			name: "UpdateMeHandler: 422",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			setup: func(a args, f fields, c *gin.Context) {
				body_string := "ascsd"

				body, _ := json.Marshal(body_string)
				c.Request = httptest.NewRequest("PUT", "/update", bytes.NewBuffer(body))
			},
			expectedStatus: 422,
			expectedBody:   nil,
		},
		{
			name: "UpdateMeHandler: 401",
			args: args{
				ctx: context.Background(),
				id:  uuid.Nil,
			},
			setup: func(a args, f fields, c *gin.Context) {
				c.Request = httptest.NewRequest("PUT", "/update", nil)
			},
			expectedStatus: 401,
			expectedBody:   nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				interactor: usecase.NewMockUserInteractor(ctrl),
				presenter:  presenter.NewMockPresenter(ctrl),
			}

			h := handlers.NewUserHandlers(f.interactor, f.presenter)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			if tc.args.id != uuid.Nil {
				c.Set("user-id", tc.args.id.String())
			}

			tc.setup(tc.args, f, c)

			h.UpdateMeHandler(c)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			var userResponse *view.UserView

			body := w.Body.Bytes()
			if len(body) > 0 {
				err := json.Unmarshal(w.Body.Bytes(), &userResponse)
				if err != nil {
					t.Errorf("can't unmarshal expected body: %v", err)
				}
			} else {
				userResponse = nil
			}

			if !reflect.DeepEqual(userResponse, tc.expectedBody) {
				t.Errorf("expected body %+v, got %+v", tc.expectedBody, userResponse)
			}
		})
	}
}

func Test_userHandlers_DeleteMeHandler(t *testing.T) {
	type fields struct {
		interactor *usecase.MockUserInteractor
		presenter  *presenter.MockPresenter
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	type testCase struct {
		name           string
		args           args
		setup          func(a args, f fields, c *gin.Context)
		expectedStatus int
	}
	cases := []testCase{
		{
			name: "DeleteMeHandler: 200",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			setup: func(a args, f fields, c *gin.Context) {
				f.interactor.EXPECT().Delete(a.ctx, a.id).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "DeleteMeHandler: 500",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			setup: func(a args, f fields, c *gin.Context) {
				f.interactor.EXPECT().Delete(a.ctx, a.id).Return(fmt.Errorf("can't delete user"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "DeleteMeHandler: 404",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			setup: func(a args, f fields, c *gin.Context) {
				f.interactor.EXPECT().Delete(a.ctx, a.id).Return(sql.ErrNoRows)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "DeleteMeHandler: 401",
			args: args{
				ctx: context.Background(),
				id:  uuid.Nil,
			},
			setup: func(a args, f fields, c *gin.Context) {
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				interactor: usecase.NewMockUserInteractor(ctrl),
				presenter:  presenter.NewMockPresenter(ctrl),
			}

			h := handlers.NewUserHandlers(f.interactor, f.presenter)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			if tc.args.id != uuid.Nil {
				c.Set("user-id", tc.args.id.String())
			}

			tc.setup(tc.args, f, c)

			h.DeleteMeHandler(c)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}
		})
	}
}

func Test_userHandlers_GetByIdHandler(t *testing.T) {
	type fields struct {
		interactor *usecase.MockUserInteractor
		presenter  *presenter.MockPresenter
	}

	type args struct {
		ctx    context.Context
		id     uuid.UUID
		userId uuid.UUID
	}

	type testCase struct {
		name           string
		args           args
		setup          func(a args, f fields, c *gin.Context)
		expectedStatus int
		expectedBody   *view.UserView
	}

	cases := []testCase{
		{
			name: "GetByIdHandler: 200",
			args: args{
				ctx:    context.Background(),
				id:     uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				userId: uuid.MustParse("4aaaba4d-9d7f-45ff-8de6-37993d709522"),
			},
			setup: func(a args, f fields, c *gin.Context) {
				userDB := &entity.UserDB{
					ID:       a.id,
					Username: "John",
					Password: "password",
					Role:     "user",
				}
				userView := &view.UserView{
					Id:       a.id.String(),
					Username: "John",
					Role:     "user",
				}

				f.interactor.EXPECT().GetById(a.ctx, a.id).Return(userDB, nil)
				f.presenter.EXPECT().ToUserView(userDB).Return(userView)
				c.Set("user-id", a.userId.String())

				c.Request = httptest.NewRequest("GET", "/users/id", nil)
				c.AddParam("id", a.id.String())
			},
			expectedStatus: http.StatusOK,
			expectedBody: &view.UserView{
				Id:       "4a6e104d-9d7f-45ff-8de6-37993d709522",
				Username: "John",
				Role:     "user",
			},
		},
		{
			name: "GetByIdHandler: 500",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			setup: func(a args, f fields, c *gin.Context) {
				f.interactor.EXPECT().GetById(a.ctx, a.id).Return(nil, fmt.Errorf("can't get user"))
				c.Set("user-id", a.userId.String())

				c.Request = httptest.NewRequest("GET", "/users/id", nil)
				c.AddParam("id", a.id.String())
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
		{
			name: "GetByIdHandler: 404",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			setup: func(a args, f fields, c *gin.Context) {
				f.interactor.EXPECT().GetById(a.ctx, a.id).Return(nil, nil)
				c.Set("user-id", a.userId.String())

				c.Request = httptest.NewRequest("GET", "/users/id", nil)
				c.AddParam("id", a.id.String())
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   nil,
		},
		{
			name: "GetByIdHandler: 422",
			args: args{
				ctx: context.Background(),
				id:  uuid.Nil,
			},
			setup: func(a args, f fields, c *gin.Context) {
				c.Set("user-id", "invalid-id")
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				interactor: usecase.NewMockUserInteractor(ctrl),
				presenter:  presenter.NewMockPresenter(ctrl),
			}

			h := handlers.NewUserHandlers(f.interactor, f.presenter)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tc.setup(tc.args, f, c)

			h.GetByIdHandler(c)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			var userResponse *view.UserView
			body := w.Body.Bytes()

			if len(body) > 0 {
				err := json.Unmarshal(w.Body.Bytes(), &userResponse)
				if err != nil {
					t.Errorf("can't unmarshal expected body: %v", err)
				}
			} else {
				userResponse = nil
			}

			if !reflect.DeepEqual(userResponse, tc.expectedBody) {
				t.Errorf("expected body %+v, got %+v", tc.expectedBody, userResponse)
			}
		})
	}
}

func Test_userHandlers_GetByUsernameHandler(t *testing.T) {
	type fields struct {
		interactor *usecase.MockUserInteractor
		presenter  *presenter.MockPresenter
	}

	type args struct {
		ctx      context.Context
		username string
	}

	type testCase struct {
		name           string
		args           args
		setup          func(a args, f fields, c *gin.Context)
		expectedStatus int
		expectedBody   *view.UserView
	}

	cases := []testCase{
		{
			name: "GetByUsernameHandler: 200",
			args: args{
				ctx:      context.Background(),
				username: "testUser",
			},
			setup: func(a args, f fields, c *gin.Context) {
				userDB := &entity.UserDB{
					ID:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
					Username: "testUser",
					Password: "password",
					Role:     "user",
				}

				userView := &view.UserView{
					Id:       "4a6e104d-9d7f-45ff-8de6-37993d709522",
					Username: "testUser",
					Role:     "user",
				}

				f.interactor.EXPECT().GetByUsername(a.ctx, a.username).Return(userDB, nil)
				f.presenter.EXPECT().ToUserView(userDB).Return(userView)
			},
			expectedStatus: http.StatusOK,
			expectedBody: &view.UserView{
				Id:       "4a6e104d-9d7f-45ff-8de6-37993d709522",
				Username: "testUser",
				Role:     "user",
			},
		},
		{
			name: "GetByUsernameHandler: 500",
			args: args{
				ctx:      context.Background(),
				username: "nonexistentUser",
			},
			setup: func(a args, f fields, c *gin.Context) {
				f.interactor.EXPECT().GetByUsername(a.ctx, a.username).Return(nil, fmt.Errorf("can't get user"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
		{
			name: "GetByUsernameHandler: 404",
			args: args{
				ctx:      context.Background(),
				username: "nonexistentUser",
			},
			setup: func(a args, f fields, c *gin.Context) {
				f.interactor.EXPECT().GetByUsername(a.ctx, a.username).Return(nil, nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				interactor: usecase.NewMockUserInteractor(ctrl),
				presenter:  presenter.NewMockPresenter(ctrl),
			}

			h := handlers.NewUserHandlers(f.interactor, f.presenter)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tc.setup(tc.args, f, c)

			c.Params = append(c.Params, gin.Param{Key: "username", Value: tc.args.username})

			h.GetByUsernameHandler(c)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			var userResponse *view.UserView
			body := w.Body.Bytes()

			if len(body) > 0 {
				err := json.Unmarshal(w.Body.Bytes(), &userResponse)
				if err != nil {
					t.Errorf("can't unmarshal expected body: %v", err)
				}
			} else {
				userResponse = nil
			}

			if !reflect.DeepEqual(userResponse, tc.expectedBody) {
				t.Errorf("expected body %+v, got %+v", tc.expectedBody, userResponse)
			}
		})
	}
}

func Test_userHandlers_UpdateHandler(t *testing.T) {
	type fields struct {
		interactor *usecase.MockUserInteractor
		presenter  *presenter.MockPresenter
	}

	type args struct {
		ctx context.Context
		id  uuid.UUID
	}

	type testCase struct {
		name           string
		args           args
		setup          func(a args, f fields, c *gin.Context)
		expectedStatus int
		expectedBody   *view.UserView
	}

	cases := []testCase{
		{
			name: "UpdateHandler: 200",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			setup: func(a args, f fields, c *gin.Context) {
				userCreate := &entity.UserCreate{
					Username: "UpdatedUser",
					Password: "updatedPassword",
				}

				userDB := &entity.UserDB{
					ID:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
					Username: "UpdatedUser",
					Password: "updatedPassword",
					Role:     "user",
				}

				userView := &view.UserView{
					Id:       "4a6e104d-9d7f-45ff-8de6-37993d709522",
					Username: "UpdatedUser",
					Role:     "user",
				}

				f.interactor.EXPECT().Update(a.ctx, a.id, userCreate).Return(userDB, nil)
				f.presenter.EXPECT().ToUserView(userDB).Return(userView)

				body, _ := json.Marshal(userCreate)
				c.Request = httptest.NewRequest("PUT", "/update", bytes.NewBuffer(body))
				c.AddParam("id", "4a6e104d-9d7f-45ff-8de6-37993d709522")
			},
			expectedStatus: http.StatusOK,
			expectedBody: &view.UserView{
				Id:       "4a6e104d-9d7f-45ff-8de6-37993d709522",
				Username: "UpdatedUser",
				Role:     "user",
			},
		},
		{
			name: "UpdateHandler: 500",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			setup: func(a args, f fields, c *gin.Context) {
				userCreate := &entity.UserCreate{
					Username: "UpdatedUser",
					Password: "updatedPassword",
				}

				f.interactor.EXPECT().Update(a.ctx, a.id, userCreate).Return(nil, fmt.Errorf("can't update user"))
				body, _ := json.Marshal(userCreate)
				c.Request = httptest.NewRequest("PUT", "/update", bytes.NewBuffer(body))
				c.AddParam("id", "4a6e104d-9d7f-45ff-8de6-37993d709522")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
		{
			name: "UpdateHandler: 422 (Unprocessable Entity)",
			args: args{
				ctx: context.Background(),
				id:  uuid.Nil,
			},
			setup: func(a args, f fields, c *gin.Context) {
				userCreate := &entity.UserCreate{
					Username: "UpdatedUser",
					Password: "updatedPassword",
				}

				body, _ := json.Marshal(userCreate)
				c.Request = httptest.NewRequest("PUT", "/update", bytes.NewBuffer(body))
				c.AddParam("id", "invalid-uuid")
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   nil,
		},
		{
			name: "UpdateHandler: 422 (Unprocessable Entity)",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			setup: func(a args, f fields, c *gin.Context) {
				invalidBody := []byte("invalid-json-body")
				c.Request = httptest.NewRequest("PUT", "/update", bytes.NewBuffer(invalidBody))
				c.AddParam("id", "4a6e104d-9d7f-45ff-8de6-37993d709522")
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				interactor: usecase.NewMockUserInteractor(ctrl),
				presenter:  presenter.NewMockPresenter(ctrl),
			}

			h := handlers.NewUserHandlers(f.interactor, f.presenter)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tc.setup(tc.args, f, c)

			h.UpdateHandler(c)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			var userResponse *view.UserView
			body := w.Body.Bytes()

			if len(body) > 0 {
				err := json.Unmarshal(w.Body.Bytes(), &userResponse)
				if err != nil {
					t.Errorf("can't unmarshal expected body: %v", err)
				}
			} else {
				userResponse = nil
			}

			if !reflect.DeepEqual(userResponse, tc.expectedBody) {
				t.Errorf("expected body %+v, got %+v", tc.expectedBody, userResponse)
			}
		})
	}
}

func Test_userHandlers_DeleteHandler(t *testing.T) {
	type fields struct {
		interactor *usecase.MockUserInteractor
		presenter  *presenter.MockPresenter
	}

	type args struct {
		ctx context.Context
		id  uuid.UUID
	}

	type testCase struct {
		name           string
		args           args
		setup          func(a args, f fields, c *gin.Context)
		expectedStatus int
	}

	cases := []testCase{
		{
			name: "DeleteHandler: 204",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			setup: func(a args, f fields, c *gin.Context) {
				f.interactor.EXPECT().Delete(a.ctx, a.id).Return(nil)
				c.AddParam("id", a.id.String())
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "DeleteHandler: 500",
			args: args{
				ctx: context.Background(),
				id:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			setup: func(a args, f fields, c *gin.Context) {
				f.interactor.EXPECT().Delete(a.ctx, a.id).Return(fmt.Errorf("can't delete user"))
				c.AddParam("id", a.id.String())
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "DeleteHandler: 422 (Unprocessable Entity)",
			args: args{
				ctx: context.Background(),
				id:  uuid.Nil,
			},
			setup: func(a args, f fields, c *gin.Context) {
				c.AddParam("id", "invalid-id")
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				interactor: usecase.NewMockUserInteractor(ctrl),
				presenter:  presenter.NewMockPresenter(ctrl),
			}

			h := handlers.NewUserHandlers(f.interactor, f.presenter)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tc.setup(tc.args, f, c)

			h.DeleteHandler(c)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}
		})
	}
}

func Test_userHandlers_LikeTrack(t *testing.T) {
	type fields struct {
		interactor *usecase.MockUserInteractor
		presenter  *presenter.MockPresenter
	}

	type args struct {
		ctx             context.Context
		userID, trackID uuid.UUID
	}

	type testCase struct {
		name           string
		args           args
		setup          func(a args, f fields, c *gin.Context)
		expectedStatus int
	}

	cases := []testCase{
		{
			name: "LikeTrack: 200",
			args: args{
				ctx:     context.Background(),
				userID:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				trackID: uuid.MustParse("8a9c1c8b-bc2f-40c6-8df1-04e0cc75ff85"),
			},
			setup: func(a args, f fields, c *gin.Context) {
				f.interactor.EXPECT().LikeTrack(a.ctx, a.userID, a.trackID).Return(nil)
				c.Set("user-id", a.userID.String())
				c.AddParam("id", a.trackID.String())
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "LikeTrack: 500",
			args: args{
				ctx:     context.Background(),
				userID:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				trackID: uuid.MustParse("8a9c1c8b-bc2f-40c6-8df1-04e0cc75ff85"),
			},
			setup: func(a args, f fields, c *gin.Context) {
				f.interactor.EXPECT().LikeTrack(a.ctx, a.userID, a.trackID).Return(fmt.Errorf("can't like track"))
				c.Set("user-id", a.userID.String())
				c.AddParam("id", a.trackID.String())
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "LikeTrack: 422 (Unprocessable Entity)",
			args: args{
				ctx:     context.Background(),
				userID:  uuid.Nil,
				trackID: uuid.Nil,
			},
			setup: func(a args, f fields, c *gin.Context) {
				c.Set("user-id", "invalid-user-id")
				c.AddParam("id", "invalid-track-id")
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				interactor: usecase.NewMockUserInteractor(ctrl),
				presenter:  presenter.NewMockPresenter(ctrl),
			}

			h := handlers.NewUserHandlers(f.interactor, f.presenter)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tc.setup(tc.args, f, c)

			h.LikeTrack(c)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}
		})
	}
}

func Test_userHandlers_DislikeTrack(t *testing.T) {
	type fields struct {
		interactor *usecase.MockUserInteractor
		presenter  *presenter.MockPresenter
	}

	type args struct {
		ctx             context.Context
		userID, trackID uuid.UUID
	}

	type testCase struct {
		name           string
		args           args
		setup          func(a args, f fields, c *gin.Context)
		expectedStatus int
	}

	cases := []testCase{
		{
			name: "DislikeTrack: 200",
			args: args{
				ctx:     context.Background(),
				userID:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				trackID: uuid.MustParse("8a9c1c8b-bc2f-40c6-8df1-04e0cc75ff85"),
			},
			setup: func(a args, f fields, c *gin.Context) {
				f.interactor.EXPECT().DislikeTrack(a.ctx, a.userID, a.trackID).Return(nil)
				c.Set("user-id", a.userID.String())
				c.AddParam("id", a.trackID.String())
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "DislikeTrack: 500",
			args: args{
				ctx:     context.Background(),
				userID:  uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
				trackID: uuid.MustParse("8a9c1c8b-bc2f-40c6-8df1-04e0cc75ff85"),
			},
			setup: func(a args, f fields, c *gin.Context) {
				f.interactor.EXPECT().DislikeTrack(a.ctx, a.userID, a.trackID).Return(fmt.Errorf("can't dislike track"))
				c.Set("user-id", a.userID.String())
				c.AddParam("id", a.trackID.String())
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "DislikeTrack: 422 (Unprocessable Entity)",
			args: args{
				ctx:     context.Background(),
				userID:  uuid.Nil,
				trackID: uuid.Nil,
			},
			setup: func(a args, f fields, c *gin.Context) {
				c.Set("user-id", "invalid-user-id")
				c.AddParam("id", "invalid-track-id")
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				interactor: usecase.NewMockUserInteractor(ctrl),
				presenter:  presenter.NewMockPresenter(ctrl),
			}

			h := handlers.NewUserHandlers(f.interactor, f.presenter)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tc.setup(tc.args, f, c)

			h.DislikeTrack(c)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}
		})
	}
}

func Test_userHandlers_ShowLikedTracks(t *testing.T) {
	type fields struct {
		interactor *usecase.MockUserInteractor
		presenter  *presenter.MockPresenter
	}

	type args struct {
		ctx    context.Context
		userID uuid.UUID
	}

	type testCase struct {
		name           string
		args           args
		setup          func(a args, f fields, c *gin.Context)
		expectedStatus int
		expectedBody   interface{} // Update the type based on your expected response
	}

	cases := []testCase{
		{
			name: "ShowLikedTracks: 200",
			args: args{
				ctx:    context.Background(),
				userID: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			setup: func(a args, f fields, c *gin.Context) {
				likedTracks := []*entity.MusicDB{
					{
						Id:       uuid.MustParse("8a9c1c8b-bc2f-40c6-8df1-04e0cc75ff85"),
						Name:     "Track 1",
						Release:  time.Now(),
						FileName: "sample.mp3",
						Size:     1024,
						Duration: "00:01:00",
					},
					{
						Id:       uuid.MustParse("3a4c3e8b-1234-4321-5678-9abcdeffedcb"),
						Name:     "Track 2",
						Release:  time.Now(),
						FileName: "sample.mp3",
						Size:     1024,
						Duration: "00:01:00",
					},
				}
				f.interactor.EXPECT().ShowLikedTracks(a.ctx, a.userID).Return(likedTracks, nil)
				f.presenter.EXPECT().ToListMusicView(likedTracks).Return([]*view.MusicView{
					{
						ID:       "8a9c1c8b-bc2f-40c6-8df1-04e0cc75ff85",
						Name:     "Track 1",
						Size:     "1.00 KB",
						Duration: "00:01:00",
					},
					{
						ID:       "3a4c3e8b-1234-4321-5678-9abcdeffedcb",
						Name:     "Track 2",
						Size:     "1.00 KB",
						Duration: "00:01:00",
					},
				})
				c.Set("user-id", a.userID.String())
			},
			expectedStatus: http.StatusOK,
			expectedBody: []*view.MusicView{
				{
					ID:       "8a9c1c8b-bc2f-40c6-8df1-04e0cc75ff85",
					Name:     "Track 1",
					Size:     "1.00 KB",
					Duration: "00:01:00",
				},
				{
					ID:       "3a4c3e8b-1234-4321-5678-9abcdeffedcb",
					Name:     "Track 2",
					Size:     "1.00 KB",
					Duration: "00:01:00",
				},
			},
		},
		{
			name: "ShowLikedTracks: 500",
			args: args{
				ctx:    context.Background(),
				userID: uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
			},
			setup: func(a args, f fields, c *gin.Context) {
				f.interactor.EXPECT().ShowLikedTracks(a.ctx, a.userID).Return(nil, fmt.Errorf("can't show liked tracks"))
				c.Set("user-id", a.userID.String())
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
		{
			name: "ShowLikedTracks: 401",
			args: args{
				ctx:    context.Background(),
				userID: uuid.Nil,
			},
			setup: func(a args, f fields, c *gin.Context) {
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				interactor: usecase.NewMockUserInteractor(ctrl),
				presenter:  presenter.NewMockPresenter(ctrl),
			}

			h := handlers.NewUserHandlers(f.interactor, f.presenter)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tc.setup(tc.args, f, c)

			h.ShowLikedTracks(c)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			if tc.expectedBody != nil {
				// Perform assertions on the response body
				// You may need to adjust this based on your actual response structure
				var responseBody []*view.MusicView // Update the type based on your expected response
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				if err != nil {
					t.Errorf("can't unmarshal response body: %v", err)
				}

				if !reflect.DeepEqual(responseBody, tc.expectedBody) {
					t.Errorf("expected body %+v, got %+v", tc.expectedBody, responseBody)
				}
			}
		})
	}
}
