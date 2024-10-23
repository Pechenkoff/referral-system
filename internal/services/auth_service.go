package services

import (
	"errors"
	"referral-system/internal/entities"
	"referral-system/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

// AuthService интерфейс для аутентификации и регистрации
type AuthService interface {
	RegisterUser(name, email, password string) (*entities.User, error)
	LoginUser(email, password string) (*entities.User, error)
}

// authService реализация AuthService
type authService struct {
	userRepo repositories.UserRepository
}

// NewAuthService создает новый AuthService
func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// RegisterUser регистрирует нового пользователя
func (s *authService) RegisterUser(name, email, password string) (*entities.User, error) {
	// Проверим, существует ли пользователь с таким email
	_, err := s.userRepo.GetUserByEmail(email)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создаем нового пользователя
	user := &entities.User{
		Name:           name,
		Email:          email,
		HashedPassword: string(hashedPassword),
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// LoginUser проверяет учетные данные пользователя и возвращает пользователя
func (s *authService) LoginUser(email, password string) (*entities.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Проверим пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
