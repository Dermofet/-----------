package presenter

import (
	"fmt"
	"music-backend-test/internal/api/http/view"
	"music-backend-test/internal/entity"
)

type presenter struct{}

func NewPresenter() *presenter {
	return &presenter{}
}

func (p *presenter) ToMusicView(music *entity.MusicDB) *view.MusicView {
	return &view.MusicView{
		ID:       music.Id.String(),
		Name:     music.Name,
		Size:     p.formatBytes(music.Size),
		Duration: music.Duration,
	}
}

func (p *presenter) formatBytes(bytes uint64) string {
	const (
		KB = 1 << 10
		MB = 1 << 20
		GB = 1 << 30
		TB = 1 << 40
		PB = 1 << 50
		EB = 1 << 60
	)

	switch {
	case bytes >= EB:
		return fmt.Sprintf("%.2f EB", float64(bytes)/float64(EB))
	case bytes >= PB:
		return fmt.Sprintf("%.2f PB", float64(bytes)/float64(PB))
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(TB))
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func (p *presenter) ToListMusicView(musics []*entity.MusicDB) []*view.MusicView {
	view := make([]*view.MusicView, len(musics))
	for i, music := range musics {
		view[i] = p.ToMusicView(music)
	}
	return view
}

func (p *presenter) ToTokenView(token *entity.Token) (*view.TokenView, error) {
	token_string, err := token.String()
	if err != nil {
		return nil, err
	}

	return &view.TokenView{
		Token: token_string,
	}, nil
}

func (p *presenter) ToUserView(user *entity.UserDB) *view.UserView {
	return &view.UserView{
		Id:       user.ID.String(),
		Username: user.Username,
		Role:     user.Role,
	}
}

func (p *presenter) ToListUserView(users []*entity.UserDB) []*view.UserView {
	view := make([]*view.UserView, len(users))
	for i, user := range users {
		view[i] = p.ToUserView(user)
	}
	return view
}
