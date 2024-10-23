package entities

import "time"

// ReferralCode - структура для реферального кода
type ReferralCode struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"` // Ссылка на пользователя, который создал код
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"` // Срок истечения кода
}

// Referral - структура для связи между реферером и рефералом
type Referral struct {
	ID         int `json:"id"`
	ReferrerID int `json:"referrer_id"` // ID реферера
	RefereeID  int `json:"referee_id"`  // ID реферала
}
