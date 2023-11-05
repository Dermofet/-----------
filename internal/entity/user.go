package entity

import (
	"github.com/google/uuid"
)

// Представление пользователя в бд
type UserDB struct {
	ID       uuid.UUID `db:"id"`       // ID
	Username string    `db:"username"` // Имя пользователя
	Password string    `db:"password"` // Пароль
	Role     string    `db:"role"`     // Роль
}

// Представление пользователя для создания записи в бд
type UserCreate struct {
	Username string `json:"username"` // Имя пользователя
	Password string `json:"password"` // Пароль
}
