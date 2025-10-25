package middleware

import (
	"context"
	"net/http"
	"strconv"
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

// RequireUserAccessOrAdmin middleware, проверяющий что пользователь имеет доступ к данным или является админом
// Позволяет:
// - Пользователю получать свои собственные данные
// - Администратору получать данные всех пользователей
func (m *AuthMiddleware) RequireUserAccessOrAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем пользователя из контекста
		currentUser := GetUserFromContext(r.Context())
		if currentUser == nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		// Извлекаем ID пользователя из URL параметров
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "ID parameter is required", http.StatusBadRequest)
			return
		}

		targetUserID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
			return
		}

		// Проверяем права доступа
		// Админ может получать данные всех пользователей
		// Обычный пользователь может получать только свои данные
		if !currentUser.IsAdmin && currentUser.UserID != targetUserID {
			http.Error(w, "Access denied: you can only access your own data", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
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
