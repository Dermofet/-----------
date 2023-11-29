package utils

import (
	"music-backend-test/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetSupportedFileType(t *testing.T) {
	type args struct {
		filename string
	}

	tests := []struct {
		name    string
		args    *args
		want    utils.FileType
		wantErr bool
	}{
		{
			name: "Get mp3 file type",
			args: &args{
				filename: "test.mp3",
			},
			want:    "MP3",
			wantErr: false,
		},
		{
			name: "Get mp3 file type",
			args: &args{
				filename: "test.MP3",
			},
			want:    "MP3",
			wantErr: false,
		},
		{
			name: "Get Ogg file type",
			args: &args{
				filename: "test.Ogg",
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utils := utils.NewmusicUtils()
			got, gotErr := utils.GetSupportedFileType(tt.args.filename)

			if tt.wantErr == true {
				assert.Error(t, gotErr)
			} else {
				if assert.NoError(t, gotErr) {
					assert.Equal(t, tt.want, got)
				}
			}
		})
	}
}

func Test_GetAudioDuration(t *testing.T) {
	type args struct {
		fileType utils.FileType
		filepath string
		os       utils.FileSystem
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Get duration",
			args: args{
				fileType: "MP3",
				filepath: "test.mp3",
				os:       utils.NewMockOS(),
			},
			want:    "00:05:30",
			wantErr: false,
		},
		{
			name: "Error in getMP3Duration",
			args: args{
				fileType: "MP3",
				filepath: "test.mp3",
				os:       utils.NewMockOS(),
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utils := utils.NewmusicUtils()

			got, gotErr := utils.GetAudioDuration(tt.args.fileType, tt.args.filepath, tt.args.os)
			if tt.wantErr == true {
				assert.Error(t, gotErr)
			} else {
				if assert.NoError(t, gotErr) {
					assert.Equal(t, tt.want, got)
				}
			}
		})
	}
}
