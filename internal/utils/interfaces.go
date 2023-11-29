package utils

import (
	"io"
	"os"
)

type MusicUtils interface {
	GetSupportedFileType(filename string) (FileType, error)
	GetAudioDuration(fileType FileType, filePath string, os FileSystem) (string, error)
}

type FileSystem interface {
	Open(string) (*os.File, error)
	Create(string) (*os.File, error)
	Copy(writer io.Writer, reader io.Reader) (int64, error)
	Remove(string) error
}
