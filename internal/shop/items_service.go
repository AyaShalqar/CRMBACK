package shop

import (
	"context"
	"fmt"
	"time"
)

type ItemService struct {
	repo *Repository
}

func NewItemService(repo *Repository) *ItemService {
	return &ItemService{repo: repo}
}

func (s *ItemService) CreateItem(ctx context.Context, item *Item) error {

	if item.Name == "" {
		return fmt.Errorf("название товара не может быть пустым")
	}

	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	return s.repo.CreateItem(ctx, item)
}

func (s *ItemService) GetItemsForShop(ctx context.Context, shopID int) ([]Item, error) {
	// Тоже можно проверить права.
	return s.repo.GetItems(ctx, shopID)
}

func (s *ItemService) GetItemByID(ctx context.Context, itemID int) (*Item, error) {
	return s.repo.GetItemByID(ctx, itemID)
}

func (s *ItemService) UpdateItem(ctx context.Context, item Item) error {

	item.UpdatedAt = time.Now()
	return s.repo.UpdateItem(ctx, item)
}

func (s *ItemService) DeleteItem(ctx context.Context, itemID int) error {
	return s.repo.DeleteItem(ctx, itemID)
}
