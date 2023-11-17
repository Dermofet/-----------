package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"music-backend-test/internal/api/http/handlers"
	"music-backend-test/internal/api/http/presenter"
	"music-backend-test/internal/api/http/view"
	"music-backend-test/internal/entity"
	"music-backend-test/internal/usecase"
	"net/http/httptest"
	reflect "reflect"
	"testing"

	"github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func MustMadeString(token *entity.Token) string {
	t, err := token.String()
	if err != nil {
		panic(err)
	}
	return t
}

func Test_authHandlers_SignUp(t *testing.T) {
	type fields struct {
		interactor *usecase.MockUserInteractor
		presenter  *presenter.MockPresenter
	}
	type args struct {
		ctx  context.Context
		user *entity.UserCreate
	}
	type testCase struct {
		name           string
		args           args
		setup          func(a args, f fields, c *gin.Context)
		expectedStatus int
		expectedBody   *view.TokenView
	}
	cases := []testCase{
		{
			name: "SignUp: 201",
			args: args{
				ctx: context.Background(),
				user: &entity.UserCreate{
					Username: "John",
					Password: "password",
				},
			},
			setup: func(a args, f fields, c *gin.Context) {
				userId := uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522")
				token := entity.GenerateToken(userId)

				tokenView := &view.TokenView{
					Token: MustMadeString(token),
				}

				fmt.Println(a.user)

				f.interactor.EXPECT().GetByUsername(a.ctx, a.user.Username).Return(nil, nil)
				f.interactor.EXPECT().Create(a.ctx, a.user).Return(userId, nil)
				f.presenter.EXPECT().ToTokenView(token).Return(tokenView, nil)

				body, _ := json.Marshal(a.user)
				c.Request = httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
			},
			expectedStatus: 201,
			expectedBody: &view.TokenView{
				Token: MustMadeString(entity.GenerateToken(uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"))),
			},
		},
		{
			name: "SignUp: 409",
			args: args{
				ctx: context.Background(),
				user: &entity.UserCreate{
					Username: "John",
					Password: "password",
				},
			},
			setup: func(a args, f fields, c *gin.Context) {
				newUser := &entity.UserCreate{
					Username: "John",
					Password: "password",
				}

				f.interactor.EXPECT().GetByUsername(a.ctx, newUser.Username).Return(&entity.UserDB{}, nil)

				body, _ := json.Marshal(newUser)
				c.Request = httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
			},
			expectedStatus: 409,
			expectedBody:   nil,
		},
		{
			name: "SignUp: 500",
			args: args{
				ctx: context.Background(),
				user: &entity.UserCreate{
					Username: "John",
					Password: "password",
				},
			},
			setup: func(a args, f fields, c *gin.Context) {
				f.interactor.EXPECT().GetByUsername(a.ctx, a.user.Username).Return(nil, fmt.Errorf("can't get user"))

				body, _ := json.Marshal(a.user)
				c.Request = httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
			},
			expectedStatus: 500,
			expectedBody:   nil,
		},
		{
			name: "SignUp: 422",
			args: args{
				ctx:  context.Background(),
				user: nil,
			},
			setup: func(a args, f fields, c *gin.Context) {
				c.Request = httptest.NewRequest("POST", "/signup", nil)
			},
			expectedStatus: 422,
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

			h := handlers.NewAuthHandlers(f.interactor, f.presenter)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tc.setup(tc.args, f, c)

			h.SignUp(c)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			if tc.expectedBody != nil {
				var responseBody view.TokenView
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				if err != nil {
					t.Errorf("can't unmarshal expected body: %v", err)
				}

				if !reflect.DeepEqual(responseBody, *tc.expectedBody) {
					t.Errorf("expected body %+v, got %+v", tc.expectedBody, &responseBody)
				}
			}
		})
	}
}

func Test_authHandlers_SignIn(t *testing.T) {
	type fields struct {
		interactor *usecase.MockUserInteractor
		presenter  *presenter.MockPresenter
	}
	type args struct {
		ctx  context.Context
		user *entity.UserCreate
	}
	type testCase struct {
		name           string
		args           args
		setup          func(a args, f fields, c *gin.Context)
		expectedStatus int
		expectedBody   *view.TokenView
	}
	cases := []testCase{
		{
			name: "SignIn: 200",
			args: args{
				ctx: context.Background(),
				user: &entity.UserCreate{
					Username: "John",
					Password: "password",
				},
			},
			setup: func(a args, f fields, c *gin.Context) {
				userDB := &entity.UserDB{
					ID:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
					Username: "John",
					Password: "password",
					Role:     "USER",
				}
				token := entity.GenerateToken(userDB.ID)

				tokenView := &view.TokenView{
					Token: MustMadeString(token),
				}

				// Исправленные вызовы методов моков
				f.interactor.EXPECT().GetByUsername(gomock.Any(), a.user.Username).Return(userDB, nil)
				f.presenter.EXPECT().ToTokenView(gomock.Any()).Return(tokenView, nil)

				body, _ := json.Marshal(a.user)
				c.Request = httptest.NewRequest("POST", "/signin", bytes.NewBuffer(body))
			},
			expectedStatus: 200,
			expectedBody: &view.TokenView{
				Token: MustMadeString(entity.GenerateToken(uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"))),
			},
		},
		{
			name: "SignIn: 404",
			args: args{
				ctx: context.Background(),
				user: &entity.UserCreate{
					Username: "NonExistentUser",
					Password: "password",
				},
			},
			setup: func(a args, f fields, c *gin.Context) {
				f.interactor.EXPECT().GetByUsername(a.ctx, a.user.Username).Return(nil, nil)

				body, _ := json.Marshal(a.user)
				c.Request = httptest.NewRequest("POST", "/signin", bytes.NewBuffer(body))
			},
			expectedStatus: 404,
			expectedBody:   nil,
		},
		{
			name: "SignIn: 401",
			args: args{
				ctx: context.Background(),
				user: &entity.UserCreate{
					Username: "John",
					Password: "wrongPassword",
				},
			},
			setup: func(a args, f fields, c *gin.Context) {
				userDB := &entity.UserDB{
					ID:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
					Username: "John",
					Password: "password",
					Role:     "USER",
				}

				f.interactor.EXPECT().GetByUsername(a.ctx, a.user.Username).Return(userDB, nil)

				body, _ := json.Marshal(a.user)
				c.Request = httptest.NewRequest("POST", "/signin", bytes.NewBuffer(body))
			},
			expectedStatus: 401,
			expectedBody:   nil,
		},
		{
			name: "SignIn: 500",
			args: args{
				ctx: context.Background(),
				user: &entity.UserCreate{
					Username: "John",
					Password: "password",
				},
			},
			setup: func(a args, f fields, c *gin.Context) {
				f.interactor.EXPECT().GetByUsername(a.ctx, a.user.Username).Return(nil, fmt.Errorf("can't get user"))

				body, _ := json.Marshal(a.user)
				c.Request = httptest.NewRequest("POST", "/signin", bytes.NewBuffer(body))
			},
			expectedStatus: 500,
			expectedBody:   nil,
		},
		{
			name: "SignIn: 422",
			args: args{
				ctx:  context.Background(),
				user: nil,
			},
			setup: func(a args, f fields, c *gin.Context) {
				c.Request = httptest.NewRequest("POST", "/signin", nil)
			},
			expectedStatus: 422,
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

			h := handlers.NewAuthHandlers(f.interactor, f.presenter)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			tc.setup(tc.args, f, c)

			h.SignIn(c)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			if tc.expectedBody != nil {
				var responseBody view.TokenView
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				if err != nil {
					t.Errorf("can't unmarshal expected body: %v", err)
				}

				if !reflect.DeepEqual(responseBody, *tc.expectedBody) {
					t.Errorf("expected body %+v, got %+v", tc.expectedBody, responseBody)
				}
			}
		})
	}
}
