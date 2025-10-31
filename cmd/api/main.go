package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"meawle/internal/config"
	"meawle/internal/database"
	"meawle/internal/handlers"
	"meawle/internal/middleware"
	"meawle/internal/repositories"
	"meawle/internal/services"
)

func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Настройка логгера
	logger := log.New(os.Stdout, "[API] ", log.LstdFlags|log.Lshortfile)

	// Инициализация базы данных
	db, err := database.New(cfg.DBPath)
	if err != nil {
		logger.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Запуск миграций
	if err := db.RunMigrations("migrations"); err != nil {
		logger.Fatal("Failed to run migrations:", err)
	}

	// Инициализация репозиториев
	userRepo := repositories.NewUserRepository(db)
	catBreedRepo := repositories.NewCatBreedRepository(db)

	// Инициализация сервисов
	userService := services.NewUserService(userRepo, cfg.JWTSecret)
	catBreedService := services.NewCatBreedService(catBreedRepo)

	// Инициализация хэндлеров
	userHandler := handlers.NewUserHandler(userService)
	catBreedHandler := handlers.NewCatBreedHandler(catBreedService)

	// Инициализация middleware
	authMiddleware := middleware.NewAuthMiddleware(userService)
	catBreedMiddleware := middleware.NewCatBreedMiddleware(catBreedService)

	// Создание маршрутизатора
	router := setupRoutes(userHandler, catBreedHandler, authMiddleware, catBreedMiddleware)

	// Настройка сервера
	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запуск сервера в горутине
	go func() {
		logger.Printf("Server starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start:", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", err)
	}

	logger.Println("Server exited")
}

func setupRoutes(userHandler *handlers.UserHandler, catBreedHandler *handlers.CatBreedHandler,
	authMiddleware *middleware.AuthMiddleware, catBreedMiddleware *middleware.CatBreedMiddleware) http.Handler {

	mux := http.NewServeMux()

	// Публичные маршруты
	mux.HandleFunc("/api/register", userHandler.Register)
	mux.HandleFunc("/api/login", userHandler.Login)
	mux.HandleFunc("/api/cat-breeds", catBreedHandler.GetAllCatBreeds)
	mux.HandleFunc("/api/cat-breed", catBreedHandler.GetCatBreed)

	// Защищенные маршруты пользователей
	mux.Handle("/api/users", authMiddleware.RequireAuth(http.HandlerFunc(userHandler.GetAllUsers)))
	mux.Handle("/api/user",
		authMiddleware.RequireAuth(
			authMiddleware.RequireUserAccessOrAdmin(
				http.HandlerFunc(userHandler.GetUser),
			),
		),
	)
	mux.Handle("/api/user/update",
		authMiddleware.RequireAuth(
			authMiddleware.RequireUserAccessOrAdmin(
				http.HandlerFunc(userHandler.UpdateUser),
			),
		),
	)
	mux.Handle("/api/user/delete",
		authMiddleware.RequireAuth(
			authMiddleware.RequireUserAccessOrAdmin(
				http.HandlerFunc(userHandler.DeleteUser),
			),
		),
	)

	// Защищенные маршруты пород кошек
	mux.Handle("/api/cat-breed/create", authMiddleware.RequireAuth(http.HandlerFunc(catBreedHandler.Create)))
	mux.Handle("/api/cat-breed/update",
		authMiddleware.RequireAuth(
			catBreedMiddleware.RequireCatBreedOwnerOrAdmin(
				http.HandlerFunc(catBreedHandler.UpdateCatBreed),
			),
		),
	)
	mux.Handle("/api/cat-breed/delete",
		authMiddleware.RequireAuth(
			catBreedMiddleware.RequireCatBreedOwnerOrAdmin(
				http.HandlerFunc(catBreedHandler.DeleteCatBreed),
			),
		),
	)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	return mux
}
