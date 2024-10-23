package repositories

import (
	"referral-system/internal/entities"
)

// Referral - структура для связи между реферером и рефералом
type Referral struct {
	ID         int `json:"id"`
	ReferrerID int `json:"referrer_id"` // ID реферера
	RefereeID  int `json:"referee_id"`  // ID реферала
}

// ReferralRepository интерфейс для работы с рефералами
type ReferralRepository interface {
	CreateReferralLink(referrerID, refereeID int) error
	GetReferralsByReferrerID(referrerID int) ([]*entities.Referral, error)
}
