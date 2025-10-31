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
	authMiddleware *middleware.AuthMiddleware,
) http.Handler {
	mux := http.NewServeMux()

	// Публичные маршруты
	mux.HandleFunc("/api/register", userHandler.Register)
	mux.HandleFunc("/api/login", userHandler.Login)
	mux.HandleFunc("/api/cat-breeds", catBreedHandler.GetAllCatBreeds)
	mux.HandleFunc("/api/cat-breed", catBreedHandler.GetCatBreed)

	// Защищенные маршруты пользователей
	mux.Handle("/api/users", authMiddleware.RequireAuth(http.HandlerFunc(userHandler.GetAllUsers)))
	mux.Handle("/api/user", authMiddleware.RequireAuth(http.HandlerFunc(userHandler.GetUser)))
	mux.Handle("/api/user/update", authMiddleware.RequireAuth(http.HandlerFunc(userHandler.UpdateUser)))
	mux.Handle("/api/user/delete", authMiddleware.RequireAuth(http.HandlerFunc(userHandler.DeleteUser)))
	// Защищенные маршруты пород кошек
	mux.Handle("/api/cat-breed/create", authMiddleware.RequireAuth(http.HandlerFunc(catBreedHandler.Create)))
	mux.Handle("/api/cat-breed/update", authMiddleware.RequireAuth(http.HandlerFunc(catBreedHandler.UpdateCatBreed)))
	mux.Handle("/api/cat-breed/delete", authMiddleware.RequireAuth(http.HandlerFunc(catBreedHandler.DeleteCatBreed)))
	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	return mux
}
