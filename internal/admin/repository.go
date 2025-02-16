package admin

import (
	"context"
	"crm-backend/internal/db"
	"fmt"
	"time"
)

type Repository struct {
	db *db.DB
}

func NewRepository(db *db.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Migrate() error {
	_, err := r.db.Conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			first_name VARCHAR(50),
			last_name VARCHAR(50),
			email VARCHAR(100) UNIQUE,
			password VARCHAR(100),
			role VARCHAR(20)
		);
	`)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы users: %w", err)
	}

	fmt.Println("Миграция users выполнена")
	return nil
}

func (r *Repository) InitSuperAdmin() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var exists bool

	err := r.db.Conn.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM users WHERE role = 'superadmin')
	`).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка проверки супер-админа: %w", err)
	}

	if exists {
		fmt.Println("Супер-админ уже существует")
		return nil
	}
	superAdmin := User{
		ID:        1,
		FirstName: "SuperAdmin",
		LastName:  "SuperAdmin",
		Email:     "admin@crm.kz",
		Password:  "superAdmin123",
		Role:      "superadmin",
	}
	if err := r.CreateUser(ctx, superAdmin); err != nil {
		return fmt.Errorf("ошибка создания супер-админа: %w", err)
	}
	fmt.Println("Супер-админ создан")
	return nil
}

func (r *Repository) CreateUser(ctx context.Context, user User) error {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("ошибка хэширования пароля: %w", err)
	}
	_, err = r.db.Conn.Exec(ctx, `
		INSERT INTO users (id, first_name, last_name, email, password, role)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, user.ID, user.FirstName, user.LastName, user.Email, hashedPassword, user.Role)

	if err != nil {
		return fmt.Errorf("ошибка создания пользователя: %w", err)
	}
	return nil
}
