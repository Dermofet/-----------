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
}

// Представление пользователя
type User struct {
	ID       *UserID // ID
	Username string  // Имя пользователя
}

// Представление пользователя для создания записи в бд
type UserCreate struct {
	Username string // Имя пользователя
	Password string // Пароль
}

// Представление пользователя для создания записи в бд
// type UserCreate struct {
// 	FirstName  string // Имя
// 	SecondName string // Отчество
// 	LastName   string // Фамилия
// 	Password   string // Пароль
// 	Age        int    // Возраст
// 	Email      string // Электронная почта
// 	Phone      string // Номер телефона
// }

// type UserSignIn struct {
// 	Email    string // Электронная почта
// 	Password string // Пароль
// }
