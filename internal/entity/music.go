package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CustomDate struct {
	time.Time
}

func (c *CustomDate) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	if s == "null" {
		return
	}
	c.Time, err = time.Parse("2006-01-02", s)
	return
}

func (c CustomDate) MarshalJSON() ([]byte, error) {
	if c.Time.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, c.Time.Format("2006-01-02"))), nil
}

type MusicShow struct {
	Name string `db:"name"`
}

type MusicID struct {
	Id uuid.UUID `db:"id"`
}

type MusicDB struct {
	Id      uuid.UUID  `db:"id"`
	Name    string     `db:"name"`
	Release CustomDate `db:"release"`
}

type MusicCreate struct {
	Name    string     `db:"name"`
	Release CustomDate `db:"release"`
}

func (m *MusicID) String() string {
	return m.Id.String()
}
