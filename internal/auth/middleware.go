package auth

import (
	"context"
	"net/http"
	"strings"
)

type UserContextKey struct{}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "требуется авторизация", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ParseJWT(tokenString)
		if err != nil {
			http.Error(w, "недействительный токен", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey{}, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(UserContextKey{}).(*Claims)
	return claims, ok
}
