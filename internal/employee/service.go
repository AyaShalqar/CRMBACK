package employee

import (
	"context"
	"crm-backend/internal/db"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Service содержит бизнес-логику
type Service struct {
	repo *Repository
	db   *db.DB // Нужно, чтобы создавать запись в таблице users (с ролью employee)
}

// NewService — конструктор для EmployeeService
func NewService(repo *Repository, db *db.DB) *Service {
	return &Service{repo: repo, db: db}
}

// AddEmployee — старый метод (если вы его используете).
// Но теперь нужно понимать, что "employee" здесь
// уже должен содержать UserID (созданный в таблице users).
func (s *Service) AddEmployee(ctx context.Context, ownerID int, employee Employee) error {
	isOwner, err := s.repo.IsOwner(ctx, employee.ShopID, ownerID)
	if err != nil {
		return fmt.Errorf("ошибка проверки владельца: %w", err)
	}
	if !isOwner {
		return fmt.Errorf("доступ запрещён: вы не владелец магазина")
	}

	// Тут employee.UserID, ShopID, Position
	return s.repo.AddEmployeeRecord(ctx, employee.UserID, employee.ShopID, employee.Position)
}

// GetEmployeesByShop — список сотрудников магазина (только для владельца)
func (s *Service) GetEmployeesByShop(ctx context.Context, ownerID, shopID int) ([]Employee, error) {
	isOwner, err := s.repo.IsOwner(ctx, shopID, ownerID)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки владельца: %w", err)
	}
	if !isOwner {
		return nil, fmt.Errorf("доступ запрещён: вы не владелец магазина")
	}

	return s.repo.GetEmployeesByShop(ctx, shopID)
}

func (s *Service) RemoveEmployee(ctx context.Context, ownerID, employeeID int) error {
	shopID, err := s.repo.GetShopIDByEmployee(ctx, employeeID)
	if err != nil {
		return fmt.Errorf("ошибка получения shop_id сотрудника: %w", err)
	}

	isOwner, err := s.repo.IsOwner(ctx, shopID, ownerID)
	if err != nil {
		return fmt.Errorf("ошибка проверки владельца: %w", err)
	}
	if !isOwner {
		return fmt.Errorf("доступ запрещён: вы не владелец магазина")
	}

	return s.repo.RemoveEmployee(ctx, employeeID)
}

// CreateUserForEmployee — создаёт запись в таблице users (роль employee)
func (s *Service) CreateUserForEmployee(ctx context.Context, req AddEmployeeRequest) (int, error) {
	// Хешируем пароль
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("ошибка bcrypt: %w", err)
	}

	// Вставляем в таблицу users
	var newUserID int
	err = s.db.Conn.QueryRow(ctx, `
		INSERT INTO users (first_name, last_name, email, password_hash, role)
		VALUES ($1, $2, $3, $4, 'employee')
		RETURNING id
	`, req.FirstName, req.LastName, req.Email, hashed).Scan(&newUserID)

	if err != nil {
		return 0, fmt.Errorf("ошибка вставки пользователя: %w", err)
	}

	return newUserID, nil
}

// AddEmployeeLink — вставляет запись в таблицу employees (user_id, shop_id)
func (s *Service) AddEmployeeLink(ctx context.Context, userID, shopID int, position string, ownerID int) error {
	isOwner, err := s.repo.IsOwner(ctx, shopID, ownerID)
	if err != nil {
		return err
	}
	if !isOwner {
		return fmt.Errorf("доступ запрещён: вы не владелец магазина %d", shopID)
	}

	return s.repo.AddEmployeeRecord(ctx, userID, shopID, position)
}
