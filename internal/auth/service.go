package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

	// "golang.org/x/crypto/bcrypt" // или argon2
	// JWT-библиотека (github.com/golang-jwt/jwt/v4) или другая
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type Service struct {
	repo         *Repository
	jwtSecretKey string
}

// NewService — инициализируем сервис с репо и секретным ключом JWT
func NewService(r *Repository) *Service {
	return &Service{
		repo:         r,
		jwtSecretKey: "SECRET_KEY_HERE", // Замените на ENV-переменную!
	}
}

// LoginUser — проверяет логин/пароль, возвращает JWT-токен и данные пользователя
func (s *Service) LoginUser(ctx context.Context, email, password string) (*User, string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", ErrInvalidCredentials // скрываем детали
	}

	// Проверяем пароль
	// Здесь упрощённо: sha256(plain) == user.PasswordHash?
	// В реальном проекте лучше bcrypt.CompareHashAndPassword(user.PasswordHash, password).
	passHash := hashPassword(password) // или bcrypt.CompareHashAndPassword
	if passHash != user.PasswordHash {
		return nil, "", ErrInvalidCredentials
	}

	// Генерируем JWT
	tokenString, err := s.generateJWT(user)
	if err != nil {
		return nil, "", fmt.Errorf("generateJWT error: %w", err)
	}

	return user, tokenString, nil
}

// GetUserByID — найти пользователя по ID (используется в /auth/me)
func (s *Service) GetUserByID(ctx context.Context, userID int) (*User, error) {
	return s.repo.FindByID(ctx, userID)
}

// generateJWT — создаёт JWT-токен (примитивный пример)
func (s *Service) generateJWT(u *User) (string, error) {
	// Обычная JWT-проверка. Claims, etc.
	claims := jwt.MapClaims{
		"user_id": u.ID,
		"role":    u.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24h
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecretKey))
}

// hashPassword — очень упрощённый для примера
func hashPassword(password string) string {
	// В реальном проекте лучше использовать bcrypt или argon2
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:])
}
