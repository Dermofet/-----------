package utils

import (
	"fmt"
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
		want    FileType
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
			got, gotErr := GetSupportedFileType(tt.args.filename)

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

func MockGetMP3Dutation(filepath string) (float64, error) {
	if filepath == "" {
		return 0, fmt.Errorf("filepath is nil")
	} else {
		return 330, nil
	}
}
