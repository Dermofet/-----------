package presenter

import (
	"music-backend-test/internal/api/http/view"
	"music-backend-test/internal/entity"
)

type musicPresenter struct{}

func NewMusicPresenter() *musicPresenter {
	return &musicPresenter{}
}

func (p *musicPresenter) ToMusicView(music *entity.Music) *view.MusicView {
	return &view.MusicView{
		ID:   music.Id.String(),
		Name: music.Name,
	}
}

func (p *musicPresenter) ToListMusicView(musics []*entity.Music) []*view.MusicView {
	view := make([]*view.MusicView, len(musics))
	for i, music := range musics {
		view[i] = p.ToMusicView(music)
	}
	return view
}
