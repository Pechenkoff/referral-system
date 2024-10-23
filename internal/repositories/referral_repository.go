package repositories

import (
	"referral-system/internal/entities"
)

// ReferralRepository интерфейс для работы с рефералами
type ReferralRepository interface {
	CreateReferralLink(referrerID, refereeID int) error
	GetReferralsByReferrerID(referrerID int) ([]*entities.Referral, error)
}
