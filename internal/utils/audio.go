package utils

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/tcolgate/mp3"
)

type FileType string

const (
	Invalid FileType = ""
	MP3     FileType = "MP3"
)

func GetSupportedFileType(filename string) (FileType, error) {
	// Преобразуйте расширение в нижний регистр, чтобы учесть все варианты (например, .MP3)
	extension := strings.ToTitle(filename[strings.LastIndex(filename, ".")+1:])

	// Сравните расширение с известными аудиоформатами
	audioExtensions := []FileType{MP3}

	if slices.Contains(audioExtensions, FileType(extension)) {
		return FileType(extension), nil
	}
	return Invalid, fmt.Errorf("unsupported file type: %s", extension)
}

func GetAudioDuration(fileType FileType, filePath string) (string, error) {
	switch fileType {
	case MP3:
		res, err := getMP3Duration(filePath)
		if err != nil {
			return "", fmt.Errorf("can't get audio duration: %w", err)
		}

		duration := time.Duration(res * float64(time.Second))

		return formatDuration(duration), nil
	default:
		return "", fmt.Errorf("unsupported file type: %s", fileType)
	}
}

func formatDuration(duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func getMP3Duration(filePath string) (float64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("can't open file: %w", err)
	}

	decoder := mp3.NewDecoder(file)
	if err != nil {
		return 0, fmt.Errorf("can't create decoder: %w", err)
	}

	var f mp3.Frame
	skipped := 0
	duration := 0.0

	for {
		if err := decoder.Decode(&f, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return 0, fmt.Errorf("can't decode file: %w", err)
		}

		duration = duration + f.Duration().Seconds()
	}
	return duration, nil
}
