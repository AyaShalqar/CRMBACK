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

func (s *Service) CreateUserByAdmin(ctx context.Context, dto CreateUserDto) error {
	user := User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Password:  dto.Password,
		Role:      dto.Role,
	}
	return s.repo.CreateUser(ctx, user)
}
