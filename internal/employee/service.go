package employee

import (
	"context"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// AddEmployee добавляет сотрудника в магазин (только владелец)
func (s *Service) AddEmployee(ctx context.Context, ownerID int, employee Employee) error {
	isOwner, err := s.repo.IsOwner(ctx, employee.ShopID, ownerID)
	if err != nil {
		return fmt.Errorf("ошибка проверки владельца: %w", err)
	}
	if !isOwner {
		return fmt.Errorf("доступ запрещён: вы не владелец магазина")
	}

	return s.repo.AddEmployee(ctx, employee)
}

// GetEmployeesByShop возвращает список сотрудников магазина (только владелец)
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
