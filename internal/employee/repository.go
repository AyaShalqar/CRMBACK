package employee

import (
	"context"
	"crm-backend/internal/db"
	"fmt"
)

type Repository struct {
	db *db.DB
}

func NewRepository(db *db.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Migrate() error {
	_, err := r.db.Conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS employees (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			shop_id INT REFERENCES shops(id) ON DELETE CASCADE,
			role VARCHAR(50) NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("ошибка миграции employees: %w", err)
	}
	fmt.Println("Миграция employees выполнена успешно")
	return nil
}

func (r *Repository) AddEmployee(ctx context.Context, employee Employee) error {
	query := `INSERT INTO employees (name, email, shop_id, role) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Conn.Exec(ctx, query, employee.Name, employee.Email, employee.ShopID, employee.Role)
	if err != nil {
		return fmt.Errorf("ошибка добавления сотрудника: %w", err)
	}
	return nil
}
