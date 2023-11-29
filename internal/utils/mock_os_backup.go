package utils

import (
	"fmt"
	"io"
	"os"
	"syscall"
)

type MockOS struct{}

func NewMockOS() *MockOS {
	return &MockOS{}
}

func (mockOs *MockOS) Open(file string) (*os.File, error) {
	if file == "" {
		return nil, fmt.Errorf("Error in os Open")
	}
	return os.NewFile(uintptr(syscall.Stdout), "Test.MP3"), nil
}

func (mockOs *MockOS) Copy(writer io.Writer, reader io.Reader) (int64, error) {
	if writer == nil || reader == nil {
		return 0, fmt.Errorf("Error in os Copy")
	}
	return 10, nil
}

func (mockOs *MockOS) Create(file string) (*os.File, error) {
	if file == "" {
		return nil, fmt.Errorf("Error in os Create")
	}
	return os.NewFile(uintptr(syscall.Stdout), "Test.MP3"), nil
}

func (mockOs *MockOS) Remove(file string) error {
	if file == "" {
		return fmt.Errorf("Error in os Remove")
	}
	return nil
}
