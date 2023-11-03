package entity

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type MusicID struct {
	Id uuid.UUID `db:"id"`
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

type MusicParse struct {
	Id         uuid.UUID
	Name       string
	Release    time.Time
	File       multipart.File
	FileHeader *multipart.FileHeader
}

type MusicCreate struct {
	Name     string    `db:"name"`
	Release  time.Time `db:"release_date"`
	FileName string    `db:"file_name"`
}

type MusicDB struct {
	Id       uuid.UUID `db:"id"`
	Name     string    `db:"name"`
	Release  time.Time `db:"release_date"`
	FileName string    `db:"file_name"`
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
