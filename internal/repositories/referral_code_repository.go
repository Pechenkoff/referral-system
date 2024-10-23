package repositories

import (
	"referral-system/internal/entities"
)

// ReferralCodeRepository интерфейс для работы с реферальными кодами
type ReferralCodeRepository interface {
	CreateReferralCode(referral *entities.ReferralCode) error
	GetReferralCodeByUserID(userID int) (*entities.ReferralCode, error)
	DeleteReferralCodeByUserID(userID int) error
}
