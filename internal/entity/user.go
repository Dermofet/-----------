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
	ID         uuid.UUID `db:"id"`          // ID
	FirstName  string    `db:"first_name"`  // Имя
	SecondName string    `db:"second_name"` // Отчество
	LastName   string    `db:"last_name"`   // Фамилия
	Password   string    `db:"password"`    // Пароль
	Age        int       `db:"age"`         // Возраст
	Email      string    `db:"email"`       // Электронная почта
	Phone      string    `db:"phone"`       // Номер телефона
}

type User struct {
	ID         *UserID // ID
	FirstName  string  // Имя
	SecondName string  // Отчество
	LastName   string  // Фамилия
	Password   string  // Пароль
	Age        int     // Возраст
	Email      string  // Электронная почта
	Phone      string  // Номер телефона
}

// Представление пользователя для создания записи в бд
type UserCreate struct {
	FirstName  string // Имя
	SecondName string // Отчество
	LastName   string // Фамилия
	Password   string // Пароль
	Age        int    // Возраст
	Email      string // Электронная почта
	Phone      string // Номер телефона
}

type UserSignIn struct {
	Email    string // Электронная почта
	Password string // Пароль
}
