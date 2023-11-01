// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package db is a generated GoMock package.
package db

import (
	context "context"
	entity "music-backend-test/internal/entity"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserSource is a mock of UserSource interface.
type MockUserSource struct {
	ctrl     *gomock.Controller
	recorder *MockUserSourceMockRecorder
}

// MockUserSourceMockRecorder is the mock recorder for MockUserSource.
type MockUserSourceMockRecorder struct {
	mock *MockUserSource
}

// NewMockUserSource creates a new mock instance.
func NewMockUserSource(ctrl *gomock.Controller) *MockUserSource {
	mock := &MockUserSource{ctrl: ctrl}
	mock.recorder = &MockUserSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserSource) EXPECT() *MockUserSourceMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockUserSource) CreateUser(ctx context.Context, user *entity.UserCreate) (*entity.UserID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, user)
	ret0, _ := ret[0].(*entity.UserID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserSourceMockRecorder) CreateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserSource)(nil).CreateUser), ctx, user)
}

// DeleteUser mocks base method.
func (m *MockUserSource) DeleteUser(ctx context.Context, id *entity.UserID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockUserSourceMockRecorder) DeleteUser(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockUserSource)(nil).DeleteUser), ctx, id)
}

// GetUserById mocks base method.
func (m *MockUserSource) GetUserById(ctx context.Context, id *entity.UserID) (*entity.UserDB, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserById", ctx, id)
	ret0, _ := ret[0].(*entity.UserDB)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserById indicates an expected call of GetUserById.
func (mr *MockUserSourceMockRecorder) GetUserById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserById", reflect.TypeOf((*MockUserSource)(nil).GetUserById), ctx, id)
}

// GetUserByUsername mocks base method.
func (m *MockUserSource) GetUserByUsername(ctx context.Context, email string) (*entity.UserDB, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsername", ctx, email)
	ret0, _ := ret[0].(*entity.UserDB)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsername indicates an expected call of GetUserByUsername.
func (mr *MockUserSourceMockRecorder) GetUserByUsername(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*MockUserSource)(nil).GetUserByUsername), ctx, email)
}

// UpdateUser mocks base method.
func (m *MockUserSource) UpdateUser(ctx context.Context, id *entity.UserID, user *entity.UserCreate) (*entity.UserDB, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, id, user)
	ret0, _ := ret[0].(*entity.UserDB)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserSourceMockRecorder) UpdateUser(ctx, id, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserSource)(nil).UpdateUser), ctx, id, user)
}

func (m *MockUserSource) LikeTrack(ctx context.Context, userId *entity.UserID, trackId *entity.MusicID) error{
	return nil
}
func (m *MockUserSource) DislikeTrack(ctx context.Context, userId *entity.UserID, trackId *entity.MusicID) error{
	return nil
}
func (m *MockUserSource) ShowLikedTracks(ctx context.Context, id *entity.UserID) ([]entity.MusicShow, error){
	return nil, nil
}
