package utils

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMusicUtils is a mock of MockMusicUtils interface.
type MockMusicUtils struct {
	ctrl     *gomock.Controller
	recorder *MockMusicUtilsMockRecorder
}

// MockMusicUtilsMockRecorder is the mock recorder for MockMusicUtils.
type MockMusicUtilsMockRecorder struct {
	mock *MockMusicUtils
}

// NewMockMusicUtils creates a new mock instance.
func NewMockMusicUtils(ctrl *gomock.Controller) *MockMusicUtils {
	mock := &MockMusicUtils{ctrl: ctrl}
	mock.recorder = &MockMusicUtilsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMusicUtils) EXPECT() *MockMusicUtilsMockRecorder {
	return m.recorder
}

// GetAudioDuration mocks base method.
func (m *MockMusicUtils) GetAudioDuration(fileType FileType, filePath string, filesystem FileSystem) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAudioDuration", fileType, filePath, filesystem)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of GetAudioDuration.
func (mr *MockMusicUtilsMockRecorder) GetAudioDuration(fileType FileType, filePath string, filesystem FileSystem) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAudioDuration", reflect.TypeOf((*MockMusicUtils)(nil).GetAudioDuration), fileType, filePath, filesystem)
}

// Create mocks base method.
func (m *MockMusicUtils) GetSupportedFileType(filename string) (FileType, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSupportedFileType", filename)
	ret0, _ := ret[0].(FileType)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockMusicUtilsMockRecorder) GetSupportedFileType(filename string) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSupportedFileType", reflect.TypeOf((*MockMusicUtils)(nil).GetSupportedFileType), filename)
}
