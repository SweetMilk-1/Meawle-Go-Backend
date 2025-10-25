package middleware

import (
	"net/http"
	"strconv"

	"meawle/internal/services"
)

// AccessMiddleware представляет middleware для проверки прав доступа
type AccessMiddleware struct {
	service *services.UserService
}

// NewAccessMiddleware создает новый экземпляр middleware для проверки прав доступа
func NewAccessMiddleware(service *services.UserService) *AccessMiddleware {
	return &AccessMiddleware{service: service}
}

// RequireUserAccess middleware, проверяющий что пользователь имеет доступ к данным
// Позволяет:
// - Пользователю изменять свои собственные данные
// - Администратору изменять данные всех пользователей
func (m *AccessMiddleware) RequireUserAccess(next http.Handler) http.Handler {
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
		// Админ может изменять данные всех пользователей
		// Обычный пользователь может изменять только свои данные
		if !currentUser.IsAdmin && currentUser.UserID != targetUserID {
			http.Error(w, "Access denied: you can only modify your own data", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequireUserAccessOrAdmin middleware, проверяющий что пользователь имеет доступ к данным или является админом
// Позволяет:
// - Пользователю получать свои собственные данные
// - Администратору получать данные всех пользователей
func (m *AccessMiddleware) RequireUserAccessOrAdmin(next http.Handler) http.Handler {
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
