package middleware

import (
	"net/http"
	"strconv"

	"meawle/internal/services"
)

// CatBreedMiddleware представляет middleware для проверки прав доступа к породам кошек
type CatBreedMiddleware struct {
	service *services.CatBreedService
}

// NewCatBreedMiddleware создает новый экземпляр middleware для пород кошек
func NewCatBreedMiddleware(service *services.CatBreedService) *CatBreedMiddleware {
	return &CatBreedMiddleware{service: service}
}

// RequireCatBreedOwnerOrAdmin middleware, проверяющий что пользователь является владельцем породы или администратором
func (m *CatBreedMiddleware) RequireCatBreedOwnerOrAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем пользователя из контекста
		currentUser := GetUserFromContext(r.Context())
		if currentUser == nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		// Извлекаем ID породы из URL параметров
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "ID parameter is required", http.StatusBadRequest)
			return
		}

		catBreedID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
			return
		}

		// Проверяем права доступа
		// Админ может работать с любыми породами
		// Обычный пользователь может работать только со своими породами
		canModify, err := m.service.CanUserModifyCatBreed(catBreedID, currentUser.UserID, currentUser.IsAdmin)
		if err != nil {
			http.Error(w, "Error checking access rights", http.StatusInternalServerError)
			return
		}

		if !canModify {
			http.Error(w, "Access denied: you can only modify your own cat breeds", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}