package auth

import (
	"context"
	"errors"
	"fmt"

	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type Service struct {
	repo         *Repository
	jwtSecretKey string
}

func NewService(r *Repository) *Service {
	return &Service{
		repo:         r,
		jwtSecretKey: "SECRET_KEY_HERE",
	}
}

func (s *Service) LoginUser(ctx context.Context, email, password string) (*User, string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {

			return nil, "", ErrInvalidCredentials
		}
		fmt.Printf("Ошибка сравнения пароля для пользователя %s: %v\n", email, err)

		return nil, "", fmt.Errorf("внутренняя ошибка сервера")
	}

	// Генерируем JWT
	tokenString, err := s.generateJWT(user)
	if err != nil {
		return nil, "", fmt.Errorf("generateJWT error: %w", err)
	}

	return user, tokenString, nil
}

func (s *Service) GetUserByID(ctx context.Context, userID int) (*User, error) {
	return s.repo.FindByID(ctx, userID)
}

func (s *Service) generateJWT(u *User) (string, error) {

	claims := jwt.MapClaims{
		"user_id": u.ID,
		"role":    u.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24h
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecretKey))
}
