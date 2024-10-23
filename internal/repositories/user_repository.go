package repositories

import (
	"referral-system/internal/entities"
	"time"
)

// User - структура для пользователя
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword string
	CreatedAt      time.Time
	UpdateAt       time.Time
}

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	CreateUser(user *entities.User) error
	GetUserByEmail(email string) (*entities.User, error)
	GetUserByID(id int) (*entities.User, error)
}
