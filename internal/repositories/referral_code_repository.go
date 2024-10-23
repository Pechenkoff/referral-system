package repositories

import (
	"referral-system/internal/entities"
	"time"
)

// ReferralCode - структура для реферального кода
type ReferralCode struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// ReferralCodeRepository интерфейс для работы с реферальными кодами
type ReferralCodeRepository interface {
	CreateReferralCode(referral *entities.ReferralCode) error
	GetReferralCodeByUserID(userID int) (*entities.ReferralCode, error)
	DeleteReferralCodeByUserID(userID int) error
}
