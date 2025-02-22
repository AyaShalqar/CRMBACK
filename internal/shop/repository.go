package shop

import (
	"context"
	"crm-backend/internal/db"
	"fmt"
)

type Repository struct {
	db db.DB
}

func NewRepository(db *db.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateShop(ctx context.Context, shop Shop) error {
	_, err := r.db.Conn.Exec(ctx, `
		INSERT INTO shops (name, description, owner_id)
		VALUES ($1, $2, $3)
	`, shop.Name, shop.Description, shop.OwnerID)

	if err != nil {
		return fmt.Errorf("ошибка создания магазина: %w", err)
	}
	return nil
}
