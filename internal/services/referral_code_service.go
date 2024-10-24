package services

import (
	"errors"
	"referral-system/internal/entities"
	"referral-system/internal/repositories"
	"time"

	"golang.org/x/exp/rand"
)

// referralService реализация ReferralService
type referralService struct {
	referralCodeRepo repositories.ReferralCodeRepository
	userRepo         repositories.UserRepository
	referralRepo     repositories.ReferralRepository
}

// ReferralService интерфейс для управления реферальными кодами
type ReferralService interface {
	CreateReferralCode(userID int, expiresIn time.Duration) (*entities.ReferralCode, error)
	DeleteReferralCode(userID int) error
	GetReferralCodeByUserID(userID int) (*entities.ReferralCode, error)
	RegisterWithReferralCode(referralCode string, name, email, password string) (*entities.User, error)
	GetReferralsByReferrerID(referrerID int) ([]*entities.Referral, error)
}

// NewReferralService создает новый ReferralService
func NewReferralService(referralCodeRepo repositories.ReferralCodeRepository,
	userRepo repositories.UserRepository,
	referralRepo repositories.ReferralRepository) ReferralService {
	return &referralService{
		referralRepo:     referralRepo,
		userRepo:         userRepo,
		referralCodeRepo: referralCodeRepo,
	}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateReferralCode генерирует случайную строку (реферальный код) заданной длины
func GenerateReferralCode(length int) string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(code)
}

// CreateReferralCode создает реферальный код для пользователя
func (s *referralService) CreateReferralCode(userID int, expiresIn time.Duration) (*entities.ReferralCode, error) {
	// Проверим, есть ли уже активный код
	existingCode, err := s.referralCodeRepo.GetReferralCodeByUserID(userID)
	if err != nil {
		return nil, err
	}

	if existingCode != nil {
		return nil, errors.New("referral code already exists for user")
	}

	code := GenerateReferralCode(10)

	// Создаем новый реферальный код
	expiresAt := time.Now().Add(expiresIn)
	referral := &entities.ReferralCode{
		UserID:    userID,
		Code:      code,
		ExpiresAt: expiresAt,
	}
	err = s.referralCodeRepo.CreateReferralCode(referral)
	if err != nil {
		return nil, err
	}

	return referral, nil
}

// DeleteReferralCode удаляет реферальный код пользователя
func (s *referralService) DeleteReferralCode(userID int) error {
	return s.referralCodeRepo.DeleteReferralCodeByUserID(userID)
}

// GetReferralCodeByUserID возвращает реферальный код по ID пользователя
func (s *referralService) GetReferralCodeByUserID(userID int) (*entities.ReferralCode, error) {
	referral, err := s.referralCodeRepo.GetReferralCodeByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Проверим, не истек ли срок действия кода
	if referral.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("referral code has expired")
	}

	return referral, nil
}

// RegisterWithReferralCode регистрирует нового пользователя по реферальному коду
func (s *referralService) RegisterWithReferralCode(referralCode string, name, email, password string) (*entities.User, error) {
	// Найдем реферальный код
	referral, err := s.referralCodeRepo.GetReferralByReferralCode(referralCode)
	if err != nil {
		return nil, errors.New("invalid referral code")
	}

	if referral.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("referral code has expired")
	}

	// Создаем нового пользователя
	authSvc := NewAuthService(s.userRepo, "")
	user, err := authSvc.RegisterUser(name, email, password)
	if err != nil {
		return nil, err
	}

	// Привязываем реферала к рефереру
	err = s.referralRepo.CreateReferralLink(referral.UserID, user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetReferralsByReferrerID возвращает список рефералов по ID реферера
func (s *referralService) GetReferralsByReferrerID(referrerID int) ([]*entities.Referral, error) {
	return s.referralRepo.GetReferralsByReferrerID(referrerID)
}
