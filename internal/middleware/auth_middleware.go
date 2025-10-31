package middleware

import (
	"context"
	"net/http"
	"strings"

	"meawle/internal/services"
)

// contextKey - тип для ключей контекста
type contextKey string

const (
	// UserContextKey ключ для хранения пользователя в контексте
	UserContextKey contextKey = "user"
)

// AuthMiddleware представляет middleware для аутентификации и проверки прав доступа
type AuthMiddleware struct {
	service *services.UserService
}

// NewAuthMiddleware создает новый экземпляр middleware аутентификации
func NewAuthMiddleware(service *services.UserService) *AuthMiddleware {
	return &AuthMiddleware{service: service}
}

// RequireAuth middleware, требующий аутентификации
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := m.extractToken(r)
		if token == "" {
			http.Error(w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		claims, err := m.service.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Добавляем claims в контекст
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// extractToken извлекает токен из заголовка Authorization
func (m *AuthMiddleware) extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Формат: Bearer <token>
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// GetUserFromContext извлекает пользователя из контекста
func GetUserFromContext(ctx context.Context) *services.JWTClaims {
	if user, ok := ctx.Value(UserContextKey).(*services.JWTClaims); ok {
		return user
	}
	return nil
}
