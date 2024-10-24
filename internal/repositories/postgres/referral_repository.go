package postgres

import (
	"context"
	"referral-system/internal/entities"
	"referral-system/internal/repositories"

	"github.com/jackc/pgx/v4/pgxpool"
)

// PostgresReferralRepository реализация ReferralRepository для PostgreSQL
type PostgresReferralRepository struct {
	db *pgxpool.Pool
}

// NewPostgresReferralRepository создает новый PostgresReferralRepository
func NewPostgresReferralRepository(db *pgxpool.Pool) repositories.ReferralRepository {
	return &PostgresReferralRepository{db: db}
}

// CreateReferralLink создает связь между реферером и рефералом
func (r *PostgresReferralRepository) CreateReferralLink(referrerID, refereeID int) error {
	query := `INSERT INTO referrals (referrer_id, referee_id) VALUES ($1, $2)`
	_, err := r.db.Exec(context.Background(), query, referrerID, refereeID)
	return err
}

// GetReferralsByReferrerID получает список рефералов по ID реферера
func (r *PostgresReferralRepository) GetReferralsByReferrerID(referrerID int) ([]*entities.Referral, error) {
	query := `SELECT id, referrer_id, referee_id FROM referrals WHERE referrer_id = $1`
	rows, err := r.db.Query(context.Background(), query, referrerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var referrals []*entities.Referral
	for rows.Next() {
		referral := &entities.Referral{}
		err := rows.Scan(&referral.ID, &referral.ReferrerID, &referral.RefereeID)
		if err != nil {
			return nil, err
		}
		referrals = append(referrals, referral)
	}

	return referrals, nil
}
