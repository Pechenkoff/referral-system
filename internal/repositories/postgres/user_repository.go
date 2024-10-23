package postgres

import (
	"database/sql"
	"errors"
	"referral-system/internal/entities"
	"referral-system/internal/repositories"
	"time"
)

// PostgresUserRepository реализация UserRepository для PostgreSQL
type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository создает новый PostgresUserRepository
func NewPostgresUserRepository(db *sql.DB) repositories.UserRepository {
	return &PostgresUserRepository{db: db}
}

// CreateUser создает нового пользователя в базе данных
func (r *PostgresUserRepository) CreateUser(user *entities.User) error {
	query := `INSERT INTO users (name, email, password, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRow(query, user.Name, user.Email, user.HashedPassword, time.Now(), time.Now()).Scan(&user.ID)
	return err
}

// GetUserByEmail находит пользователя по email
func (r *PostgresUserRepository) GetUserByEmail(email string) (*entities.User, error) {
	user := &entities.User{}
	query := `SELECT id, name, email, password FROM users WHERE email=$1`
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}

// GetUserByID находит пользователя по ID
func (r *PostgresUserRepository) GetUserByID(id int) (*entities.User, error) {
	user := &entities.User{}
	query := `SELECT id, name, email, password FROM users WHERE id=$1`
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}
