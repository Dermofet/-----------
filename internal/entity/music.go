package entity

import "github.com/google/uuid"

type MusicID struct {
	Id uuid.UUID
}

func (m *MusicID) String() string {
	return m.Id.String()
}

func (m *MusicID) FromString(s string) error {
	var err error
	m.Id, err = uuid.Parse(s)
	return err
}

type Music struct {
	Id   *MusicID
	Name string
}

type MusicDB struct {
	Id   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

type MusicCreate struct {
	Name string `json:"name"`
}
