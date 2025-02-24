package shop

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateShop(ctx context.Context, shop Shop) error {
	return s.repo.CreateShop(ctx, shop)
}

func (s *Service) GetShops(ctx context.Context) ([]Shop, error) {
	return s.repo.GetShops(ctx)
}

func (s *Service) UpdateShop(ctx context.Context, shop Shop) error {
	return s.repo.UpdateShop(ctx, shop)
}

func (s *Service) DeleteShop(ctx context.Context, id int) error {
	return s.repo.DeleteShop(ctx, id)
}
