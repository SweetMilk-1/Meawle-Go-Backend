package routes

import (
	"encoding/json"
	"net/http"

	"meawle/internal/handlers"
	"meawle/internal/middleware"

	"github.com/gorilla/mux"
)

// SetupRoutes возвращает настроенный HTTP роутер с gorilla/mux
func SetupRoutes(
	userHandler *handlers.UserHandler,
	catBreedHandler *handlers.CatBreedHandler,
	catHandler *handlers.CatHandler,
	authMiddleware *middleware.AuthMiddleware,
) http.Handler {
	r := mux.NewRouter()

	// API маршруты с версионированием
	api := r.PathPrefix("/api/v1").Subrouter()

	// Публичные маршруты
	api.HandleFunc("/auth/register", userHandler.Register).Methods(http.MethodPost)
	api.HandleFunc("/auth/login", userHandler.Login).Methods(http.MethodPost)
	api.HandleFunc("/users", userHandler.GetAllUsers).Methods(http.MethodGet)
	api.HandleFunc("/users/{id:[0-9]+}", userHandler.GetUser).Methods(http.MethodGet)
	api.HandleFunc("/cat-breeds", catBreedHandler.GetAllCatBreeds).Methods(http.MethodGet)
	api.HandleFunc("/cat-breeds/{id:[0-9]+}", catBreedHandler.GetCatBreed).Methods(http.MethodGet)
	api.HandleFunc("/cats", catHandler.GetAllCats).Methods(http.MethodGet)
	api.HandleFunc("/cats/{id:[0-9]+}", catHandler.GetCat).Methods(http.MethodGet)

	// Защищенные маршруты пользователей
	users := api.PathPrefix("/users").Subrouter()
	users.Use(authMiddleware.RequireAuth)
	users.HandleFunc("/{id:[0-9]+}", userHandler.UpdateUser).Methods(http.MethodPut)
	users.HandleFunc("/{id:[0-9]+}", userHandler.DeleteUser).Methods(http.MethodDelete)

	// Защищенные маршруты пород кошек
	catBreeds := api.PathPrefix("/cat-breeds").Subrouter()
	catBreeds.Use(authMiddleware.RequireAuth)
	catBreeds.HandleFunc("", catBreedHandler.Create).Methods(http.MethodPost)
	catBreeds.HandleFunc("/{id:[0-9]+}", catBreedHandler.UpdateCatBreed).Methods(http.MethodPut)
	catBreeds.HandleFunc("/{id:[0-9]+}", catBreedHandler.DeleteCatBreed).Methods(http.MethodDelete)

	// Защищенные маршруты котов
	cats := api.PathPrefix("/cats").Subrouter()
	cats.Use(authMiddleware.RequireAuth)
	cats.HandleFunc("", catHandler.Create).Methods(http.MethodPost)
	cats.HandleFunc("/user", catHandler.GetUserCats).Methods(http.MethodGet)
	cats.HandleFunc("/{id:[0-9]+}", catHandler.UpdateCat).Methods(http.MethodPut)
	cats.HandleFunc("/{id:[0-9]+}", catHandler.DeleteCat).Methods(http.MethodDelete)

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods(http.MethodGet)

	return r
}
