package repositories

import (
	"referral-system/internal/entities"
)

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	CreateUser(user *entities.User) error
	GetUserByEmail(email string) (*entities.User, error)
	GetUserByID(id int) (*entities.User, error)
}
