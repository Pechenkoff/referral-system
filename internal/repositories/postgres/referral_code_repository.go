package postgres

import (
	"context"
	"database/sql"
	"errors"
	"referral-system/internal/entities"
	"referral-system/internal/repositories"
	"time"

	"github.com/jackc/pgx/v4"
)

// PostgresReferralCodeRepository реализация ReferralRepository для PostgreSQL
type PostgresReferralCodeRepository struct {
	db *pgx.Conn
}

// NewPostgresReferralCodeRepository создает новый PostgresReferralRepository
func NewPostgresReferralCodeRepository(db *pgx.Conn) repositories.ReferralCodeRepository {
	return &PostgresReferralCodeRepository{db: db}
}

// CreateReferralCode создает новый реферальный код
func (r *PostgresReferralCodeRepository) CreateReferralCode(referral *entities.ReferralCode) error {
	query := `INSERT INTO referral_codes (user_id, code, expires_at, created_at) 
              VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRow(context.Background(), query, referral.UserID, referral.Code, referral.ExpiresAt, time.Now()).Scan(&referral.ID)
	return err
}

// GetReferralCodeByUserID получает реферальный код по ID пользователя
func (r *PostgresReferralCodeRepository) GetReferralCodeByUserID(userID int) (*entities.ReferralCode, error) {
	referral := &entities.ReferralCode{}
	query := `SELECT id, user_id, code, expires_at FROM referral_codes WHERE user_id=$1`
	err := r.db.QueryRow(context.Background(), query, userID).Scan(&referral.ID, &referral.UserID, &referral.Code, &referral.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("referral code not found")
	}
	return referral, err
}

// DeleteReferralCodeByUserID удаляет реферальный код по ID пользователя
func (r *PostgresReferralCodeRepository) DeleteReferralCodeByUserID(userID int) error {
	query := `DELETE FROM referral_codes WHERE user_id=$1`
	_, err := r.db.Exec(context.Background(), query, userID)
	return err
}

func (r *PostgresReferralCodeRepository) GetReferralByReferralCode(referralCode string) (*entities.ReferralCode, error) {
	var referral = &entities.ReferralCode{}
	query := `SELECT user_id, expires_at FROM referral_codes WHERE code=$1`
	err := r.db.QueryRow(context.Background(), query, referralCode).Scan(&referral.UserID, &referral.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	return referral, err
}
