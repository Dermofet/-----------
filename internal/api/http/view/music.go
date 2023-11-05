package view

type MusicView struct {
	ID       string `json:"id"`       // id трека
	Name     string `json:"name"`     // название трека
	Size     string `json:"size"`     // размер файла трека (в удобном для чтения виде)
	Duration string `json:"duration"` // продолжительность трека
}
