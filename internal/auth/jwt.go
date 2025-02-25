package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SecretKey = []byte("super_secret_key") // ❗ В будущем загружай из .env

type Claims struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT генерирует токен для пользователя
func GenerateJWT(id int, email, role string) (string, error) {
	claims := Claims{
		ID:    id, // ⬅️ Теперь `id` передаётся как аргумент.
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Токен живёт 24 часа
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(SecretKey)

	if err != nil {
		fmt.Println("Ошибка генерации токена:", err) // 🔥 Логируем ошибку
		return "", err
	}
	return signedToken, nil
}

// ParseJWT проверяет токен и извлекает данные
func ParseJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("недействительный токен")
	}

	return claims, nil
}
