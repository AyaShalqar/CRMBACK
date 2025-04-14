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
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			shop_id INT NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
			position VARCHAR(100),
			hired_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("ошибка миграции employees: %w", err)
	}
	fmt.Println("Миграция employees выполнена успешно")
	return nil
}

func (r *Repository) IsOwner(ctx context.Context, shopID, ownerID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM shops WHERE id = $1 AND owner_id = $2)`
	err := r.db.Conn.QueryRow(ctx, query, shopID, ownerID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ошибка проверки владельца: %w", err)
	}
	return exists, nil
}

func (r *Repository) AddEmployeeRecord(ctx context.Context, userID, shopID int, position string) error {
	query := `INSERT INTO employees (user_id, shop_id, position) VALUES ($1, $2, $3)`
	_, err := r.db.Conn.Exec(ctx, query, userID, shopID, position)
	if err != nil {
		return fmt.Errorf("ошибка добавления сотрудника в БД: %w", err)
	}
	return nil
}

func (r *Repository) GetEmployeesByShop(ctx context.Context, shopID int) ([]Employee, error) {
	query := `SELECT id, user_id, shop_id, position, hired_at FROM employees WHERE shop_id = $1`
	rows, err := r.db.Conn.Query(ctx, query, shopID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения сотрудников: %w", err)
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var emp Employee
		if err := rows.Scan(&emp.ID, &emp.UserID, &emp.ShopID, &emp.Position, &emp.HiredAt); err != nil {
			return nil, fmt.Errorf("ошибка чтения данных сотрудника: %w", err)
		}
		employees = append(employees, emp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке строк сотрудников: %w", err)
	}
	return employees, nil
}

func (r *Repository) RemoveEmployee(ctx context.Context, employeeID int) error {
	if employeeID <= 0 {
		return fmt.Errorf("неверный ID сотрудника: %d", employeeID)
	}

	query := `DELETE FROM employees WHERE id = $1`
	res, err := r.db.Conn.Exec(ctx, query, employeeID)
	if err != nil {
		return fmt.Errorf("ошибка удаления сотрудника: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("сотрудник с ID %d не найден", employeeID)
	}
	return nil
}

func (r *Repository) GetShopIDByEmployee(ctx context.Context, employeeID int) (int, error) {
	var shopID int
	query := `SELECT shop_id FROM employees WHERE id = $1`
	err := r.db.Conn.QueryRow(ctx, query, employeeID).Scan(&shopID)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения shop_id сотрудника: %w", err)
	}
	return shopID, nil
}
