package shop

import (
	"context"
	"fmt"
	"time"
)

// ItemService описывает бизнес-логику работы с товарами.
type ItemService struct {
	repo *Repository
}

// NewItemService - конструктор
func NewItemService(repo *Repository) *ItemService {
	return &ItemService{repo: repo}
}

// CreateItem - пример бизнес-логики создания товара.
func (s *ItemService) CreateItem(ctx context.Context, item *Item) error {
	// Тут можно добавить логику: проверить, что поля не пустые, нет ли слишком низкой цены и т.д.
	if item.Name == "" {
		return fmt.Errorf("название товара не может быть пустым")
	}

	// Можно проверить, есть ли магазин (shopID) и права у пользователя.
	// Но допустим, у нас это проверяется где-то в middleware или выше.

	// Установим даты, если нужно (у тебя уже NOW() в SQL, но можно и в Go)
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	return s.repo.CreateItem(ctx, item)
}

// GetItemsForShop - получить все товары конкретного магазина.
func (s *ItemService) GetItemsForShop(ctx context.Context, shopID int) ([]Item, error) {
	// Тоже можно проверить права.
	return s.repo.GetItems(ctx, shopID)
}

// GetItemByID - получить товар по ID
func (s *ItemService) GetItemByID(ctx context.Context, itemID int) (*Item, error) {
	return s.repo.GetItemByID(ctx, itemID)
}

// UpdateItem - обновить данные товара
func (s *ItemService) UpdateItem(ctx context.Context, item Item) error {
	// Можно проверить, является ли item.Name пустым,
	// и есть ли права на редактирование данного товара.
	item.UpdatedAt = time.Now()
	return s.repo.UpdateItem(ctx, item)
}

// DeleteItem - удалить товар
func (s *ItemService) DeleteItem(ctx context.Context, itemID int) error {
	return s.repo.DeleteItem(ctx, itemID)
}
