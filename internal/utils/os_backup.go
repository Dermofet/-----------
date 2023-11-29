package utils

import (
	"io"
	"os"
)

type fileSystem struct{}

func NewFileSystem() *fileSystem {
	return &fileSystem{}
}

func (fileSystem *fileSystem) Open(file string) (*os.File, error) {
	result, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (fileSystem *fileSystem) Copy(writer io.Writer, reader io.Reader) (int64, error) {
	result, err := io.Copy(writer, reader)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (fileSystem *fileSystem) Create(file string) (*os.File, error) {
	result, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (fileSystem *fileSystem) Remove(file string) error {
	err := os.Remove(file)
	if err != nil {
		return err
	}
	return nil
}
