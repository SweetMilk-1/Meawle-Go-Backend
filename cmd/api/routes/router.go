package routes

import (
	"encoding/json"
	"net/http"

	"meawle/internal/handlers"
	"meawle/internal/middleware"
)

// Router возвращает настроенный HTTP роутер
func SetupRoutes(
	userHandler *handlers.UserHandler,
	catBreedHandler *handlers.CatBreedHandler,
	catHandler *handlers.CatHandler,
	authMiddleware *middleware.AuthMiddleware,
) http.Handler {
	mux := http.NewServeMux()

	// Публичные маршруты
	mux.HandleFunc("/api/register", userHandler.Register)
	mux.HandleFunc("/api/login", userHandler.Login)
	mux.HandleFunc("/api/cat-breeds", catBreedHandler.GetAllCatBreeds)
	mux.HandleFunc("/api/cat-breed", catBreedHandler.GetCatBreed)
	mux.HandleFunc("/api/cats", catHandler.GetAllCats)
	mux.HandleFunc("/api/cat", catHandler.GetCat)

	// Защищенные маршруты пользователей
	mux.Handle("/api/users", authMiddleware.RequireAuth(http.HandlerFunc(userHandler.GetAllUsers)))
	mux.Handle("/api/user", authMiddleware.RequireAuth(http.HandlerFunc(userHandler.GetUser)))
	mux.Handle("/api/user/update", authMiddleware.RequireAuth(http.HandlerFunc(userHandler.UpdateUser)))
	mux.Handle("/api/user/delete", authMiddleware.RequireAuth(http.HandlerFunc(userHandler.DeleteUser)))

	// Защищенные маршруты пород кошек
	mux.Handle("/api/cat-breed/create", authMiddleware.RequireAuth(http.HandlerFunc(catBreedHandler.Create)))
	mux.Handle("/api/cat-breed/update", authMiddleware.RequireAuth(http.HandlerFunc(catBreedHandler.UpdateCatBreed)))
	mux.Handle("/api/cat-breed/delete", authMiddleware.RequireAuth(http.HandlerFunc(catBreedHandler.DeleteCatBreed)))

	// Защищенные маршруты котов
	mux.Handle("/api/cat/create", authMiddleware.RequireAuth(http.HandlerFunc(catHandler.Create)))
	mux.Handle("/api/cat/update", authMiddleware.RequireAuth(http.HandlerFunc(catHandler.UpdateCat)))
	mux.Handle("/api/cat/delete", authMiddleware.RequireAuth(http.HandlerFunc(catHandler.DeleteCat)))

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	return mux
}
