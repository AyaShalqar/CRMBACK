package auth

import (
	"context"
	"crm-backend/internal/db"
	"fmt"
)

// User — упрощённая модель юзера (подставьте поля, которые есть в вашей БД)
type User struct {
	ID           int    `db:"id"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
	Role         string `db:"role"` // "admin", "owner", etc.
}

// Repository — отвечает за запросы к таблице users
type Repository struct {
	db *db.DB
}

func NewRepository(database *db.DB) *Repository {
	return &Repository{db: database}
}

// FindByEmail — находит пользователя по email
func (r *Repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `
        SELECT id, email, password_hash, role
        FROM users
        WHERE email = $1
        LIMIT 1
    `
	row := r.db.Conn.QueryRow(ctx, query, email)
	var u User
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role)
	if err != nil {
		return nil, fmt.Errorf("FindByEmail: %w", err)
	}
	return &u, nil
}

// FindByID — находит пользователя по ID
func (r *Repository) FindByID(ctx context.Context, userID int) (*User, error) {
	query := `
        SELECT id, email, password_hash, role
        FROM users
        WHERE id = $1
        LIMIT 1
    `
	row := r.db.Conn.QueryRow(ctx, query, userID)
	var u User
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role)
	if err != nil {
		return nil, fmt.Errorf("FindByID: %w", err)
	}
	return &u, nil
}
