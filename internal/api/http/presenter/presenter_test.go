package presenter

import (
	"music-backend-test/internal/api/http/view"
	"music-backend-test/internal/entity"
	reflect "reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func MustParseTime(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

func Test_presenter_ToMusicView(t *testing.T) {
	type args struct {
		music *entity.MusicDB
	}
	tests := []struct {
		name string
		args args
		want *view.MusicView
	}{
		{
			name: "success ToMusicView",
			args: args{
				music: &entity.MusicDB{
					Id:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
					Name:     "Sample Music",
					Release:  MustParseTime("2006-01-02", "2022-01-01"),
					FileName: "sample.mp3",
					Size:     1024,
					Duration: "03:24",
				},
			},
			want: &view.MusicView{
				ID:       "4a6e104d-9d7f-45ff-8de6-37993d709522",
				Name:     "Sample Music",
				Size:     "1.00 KB",
				Duration: "03:24",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &presenter{}
			got := p.ToMusicView(tt.args.music)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("presenter.ToMusicView() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_presenter_ToListMusicView(t *testing.T) {
	type args struct {
		musics []*entity.MusicDB
	}
	tests := []struct {
		name string
		args args
		want []*view.MusicView
	}{
		{
			name: "success ToListMusicView",
			args: args{
				musics: []*entity.MusicDB{
					{
						Id:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
						Name:     "Sample Music 1",
						Size:     1024,
						Duration: "03:24",
					},
					{
						Id:       uuid.MustParse("6cfc3a4d-28ae-4b13-9a11-d2c335ad8658"),
						Name:     "Sample Music 2",
						Size:     2048,
						Duration: "05:12",
					},
				},
			},
			want: []*view.MusicView{
				{
					ID:       "4a6e104d-9d7f-45ff-8de6-37993d709522",
					Name:     "Sample Music 1",
					Size:     "1.00 KB",
					Duration: "03:24",
				},
				{
					ID:       "6cfc3a4d-28ae-4b13-9a11-d2c335ad8658",
					Name:     "Sample Music 2",
					Size:     "2.00 KB",
					Duration: "05:12",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &presenter{}
			got := p.ToListMusicView(tt.args.musics)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("presenter.ToListMusicView() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_presenter_ToUserView(t *testing.T) {
	type args struct {
		user *entity.UserDB
	}
	tests := []struct {
		name string
		args args
		want *view.UserView
	}{
		{
			name: "success ToUserView",
			args: args{
				user: &entity.UserDB{
					ID:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
					Username: "John",
					Password: "password",
					Role:     "admin",
				},
			},
			want: &view.UserView{
				Id:       "4a6e104d-9d7f-45ff-8de6-37993d709522",
				Username: "John",
				Role:     "admin",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &presenter{}
			got := p.ToUserView(tt.args.user)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("presenter.ToUserView() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_presenter_ToListUserView(t *testing.T) {
	type args struct {
		users []*entity.UserDB
	}
	tests := []struct {
		name string
		args args
		want []*view.UserView
	}{
		{
			name: "success ToListUserView",
			args: args{
				users: []*entity.UserDB{
					{
						ID:       uuid.MustParse("4a6e104d-9d7f-45ff-8de6-37993d709522"),
						Username: "John",
						Role:     "admin",
					},
					{
						ID:       uuid.MustParse("6cfc3a4d-28ae-4b13-9a11-d2c335ad8658"),
						Username: "Alice",
						Role:     "user",
					},
				},
			},
			want: []*view.UserView{
				{
					Id:       "4a6e104d-9d7f-45ff-8de6-37993d709522",
					Username: "John",
					Role:     "admin",
				},
				{
					Id:       "6cfc3a4d-28ae-4b13-9a11-d2c335ad8658",
					Username: "Alice",
					Role:     "user",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &presenter{}
			got := p.ToListUserView(tt.args.users)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("presenter.ToListUserView() = %v, want %v", got, tt.want)
			}
		})
	}
}
