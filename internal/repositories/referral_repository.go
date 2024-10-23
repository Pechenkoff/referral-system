package repositories

import (
	"database/sql"
	"errors"
	"referral-system/internal/entities"
	"time"
)

// ReferralCode - структура для реферального кода
type ReferralCode struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"` // Ссылка на пользователя, который создал код
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"` // Срок истечения кода
	CreatedAt time.Time `json:"created_at"`
}

// ReferralRepository интерфейс для работы с реферальными кодами
type ReferralRepository interface {
	CreateReferralCode(referral *entities.ReferralCode) error
	GetReferralCodeByUserID(userID int) (*entities.ReferralCode, error)
	DeleteReferralCodeByUserID(userID int) error
}

// PostgresReferralRepository реализация ReferralRepository для PostgreSQL
type PostgresReferralRepository struct {
	db *sql.DB
}

// NewPostgresReferralRepository создает новый PostgresReferralRepository
func NewPostgresReferralRepository(db *sql.DB) ReferralRepository {
	return &PostgresReferralRepository{db: db}
}

// CreateReferralCode создает новый реферальный код
func (r *PostgresReferralRepository) CreateReferralCode(referral *entities.ReferralCode) error {
	query := `INSERT INTO referral_codes (user_id, code, expires_at, created_at) 
              VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRow(query, referral.UserID, referral.Code, referral.ExpiresAt, time.Now()).Scan(&referral.ID)
	return err
}

// GetReferralCodeByUserID получает реферальный код по ID пользователя
func (r *PostgresReferralRepository) GetReferralCodeByUserID(userID int) (*entities.ReferralCode, error) {
	referral := &entities.ReferralCode{}
	query := `SELECT id, user_id, code, expires_at FROM referral_codes WHERE user_id=$1`
	err := r.db.QueryRow(query, userID).Scan(&referral.ID, &referral.UserID, &referral.Code, &referral.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("referral code not found")
	}
	return referral, err
}

// DeleteReferralCodeByUserID удаляет реферальный код по ID пользователя
func (r *PostgresReferralRepository) DeleteReferralCodeByUserID(userID int) error {
	query := `DELETE FROM referral_codes WHERE user_id=$1`
	_, err := r.db.Exec(query, userID)
	return err
}
