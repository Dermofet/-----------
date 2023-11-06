package entity

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

// type Music struct {
// 	Id       uuid.UUID
// 	Name     string
// 	Filename string
// 	Release  time.Time
// }

// swagger:ignore
type MusicParse struct {
	Name       string
	Release    time.Time
	File       multipart.File        `swaggerignore:"true"`
	FileHeader *multipart.FileHeader `swaggerignore:"true"`
}

type MusicDB struct {
	Id       uuid.UUID `db:"id"`           // id трека
	Name     string    `db:"name"`         // название трека
	Release  time.Time `db:"release_date"` // дата релиза трека
	FileName string    `db:"file_name"`    // имя файла
	Size     uint64    `db:"size"`         // размер файла
	Duration string    `db:"duration"`     // продолжительность трека
}

func (m *MusicDB) FilePath() string {
	return fmt.Sprintf("./internal/storage/music_storage/%s", m.FileName)
}

// type CustomDate struct {
// 	time.Time
// }
// func (c *CustomDate) UnmarshalJSON(b []byte) (err error) {
// 	s := strings.Trim(string(b), `"`)
// 	if s == "null" {
// 		return
// 	}
// 	c.Time, err = time.Parse("2006-01-02", s)
// 	return
// }

// func (c CustomDate) MarshalJSON() ([]byte, error) {
// 	if c.Time.IsZero() {
// 		return nil, nil
// 	}
// 	return []byte(fmt.Sprintf(`"%s"`, c.Time.Format("2006-01-02"))), nil
// }
