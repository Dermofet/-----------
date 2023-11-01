package view

type MusicView struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ListMusicView struct {
	Musics []*MusicView `json:"musics"`
}
