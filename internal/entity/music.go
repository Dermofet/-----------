package entity

import "github.com/google/uuid"

type MusicShow struct {
	Name string `json: "Name"`
}

type MusicID struct {
	Id uuid.UUID `json: "ID"`
}

type MusicDB struct {
	Id   uuid.UUID `json: "ID"`
	Name string    `json: "Name"`
}

type MusicCreate struct {
	Name string `json: "Name"`
}
