package postgres

import (
	"database/sql"
	"referral-system/internal/entities"
	"referral-system/internal/repositories"
)

// PostgresReferralRepository реализация ReferralRepository для PostgreSQL
type PostgresReferralRepository struct {
	db *sql.DB
}

// NewPostgresReferralRepository создает новый PostgresReferralRepository
func NewPostgresReferralRepository(db *sql.DB) repositories.ReferralRepository {
	return &PostgresReferralRepository{db: db}
}

// CreateReferralLink создает связь между реферером и рефералом
func (r *PostgresReferralRepository) CreateReferralLink(referrerID, refereeID int) error {
	query := `INSERT INTO referrals (referrer_id, referee_id) VALUES ($1, $2)`
	_, err := r.db.Exec(query, referrerID, refereeID)
	return err
}

// GetReferralsByReferrerID получает список рефералов по ID реферера
func (r *PostgresReferralRepository) GetReferralsByReferrerID(referrerID int) ([]*entities.Referral, error) {
	query := `SELECT id, referrer_id, referee_id FROM referrals WHERE referrer_id = $1`
	rows, err := r.db.Query(query, referrerID)
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
