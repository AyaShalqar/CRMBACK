package admin

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, user User) error {
	return s.repo.CreateUser(ctx, user)
}

func (s *Service) InitSuperAdmin() error {
	return s.repo.InitSuperAdmin()
}
