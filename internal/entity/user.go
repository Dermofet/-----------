package entity

import (
	"github.com/google/uuid"
)

// Представление id пользователя
type UserID struct {
	Id uuid.UUID
}

func (u *UserID) String() string {
	return u.Id.String()
}

func (u *UserID) FromString(s string) error {
	var err error
	u.Id, err = uuid.Parse(s)
	return err
}

// Представление пользователя в бд
type UserDB struct {
	ID       uuid.UUID `db:"id"`       // ID
	Username string    `db:"username"` // Имя пользователя
	Password string    `db:"password"` // Пароль
	Role     string    `db:"role"`     // Роль
}

// Представление пользователя
type User struct {
	ID       *UserID // ID
	Username string  // Имя пользователя
	Role     string  // Роль
	Password string  // Пароль
}

// Представление пользователя для создания записи в бд
type UserCreate struct {
	Username string `json:"username"` // Имя пользователя
	Password string `json:"password"` // Пароль
}
