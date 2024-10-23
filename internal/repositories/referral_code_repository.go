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

// ReferralCodeRepository интерфейс для работы с реферальными кодами
type ReferralCodeRepository interface {
	CreateReferralCode(referral *entities.ReferralCode) error
	GetReferralCodeByUserID(userID int) (*entities.ReferralCode, error)
	DeleteReferralCodeByUserID(userID int) error
}

// PostgresReferralCodeRepository реализация ReferralRepository для PostgreSQL
type PostgresReferralCodeRepository struct {
	db *sql.DB
}

// NewPostgresReferralCodeRepository создает новый PostgresReferralRepository
func NewPostgresReferralCodeRepository(db *sql.DB) ReferralCodeRepository {
	return &PostgresReferralCodeRepository{db: db}
}

// CreateReferralCode создает новый реферальный код
func (r *PostgresReferralCodeRepository) CreateReferralCode(referral *entities.ReferralCode) error {
	query := `INSERT INTO referral_codes (user_id, code, expires_at, created_at) 
              VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRow(query, referral.UserID, referral.Code, referral.ExpiresAt, time.Now()).Scan(&referral.ID)
	return err
}

// GetReferralCodeByUserID получает реферальный код по ID пользователя
func (r *PostgresReferralCodeRepository) GetReferralCodeByUserID(userID int) (*entities.ReferralCode, error) {
	referral := &entities.ReferralCode{}
	query := `SELECT id, user_id, code, expires_at FROM referral_codes WHERE user_id=$1`
	err := r.db.QueryRow(query, userID).Scan(&referral.ID, &referral.UserID, &referral.Code, &referral.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("referral code not found")
	}
	return referral, err
}

// DeleteReferralCodeByUserID удаляет реферальный код по ID пользователя
func (r *PostgresReferralCodeRepository) DeleteReferralCodeByUserID(userID int) error {
	query := `DELETE FROM referral_codes WHERE user_id=$1`
	_, err := r.db.Exec(query, userID)
	return err
}
