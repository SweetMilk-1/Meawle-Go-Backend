package main

import (
	"log"
	"meawle/internal/database"
	"meawle/internal/handlers"
	"meawle/internal/middleware"
	"meawle/internal/repositories"
	"meawle/internal/services"
	"net/http"
)

func main() {
	// Конфигурация
	dbPath := "app.db"
	jwtSecret := "your-secret-key-change-in-production"
	port := ":8080"

	// Инициализация базы данных
	db, err := database.New(dbPath)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	err = db.RunMigrations("migrations")
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Инициализация репозиториев
	userRepo := repositories.NewUserRepository(db)
	catBreedRepo := repositories.NewCatBreedRepository(db)

	// Инициализация сервисов
	userService := services.NewUserService(userRepo, jwtSecret)
	catBreedService := services.NewCatBreedService(catBreedRepo)

	// Инициализация хэндлеров
	userHandler := handlers.NewUserHandler(userService)
	catBreedHandler := handlers.NewCatBreedHandler(catBreedService)

	// Инициализация middleware
	authMiddleware := middleware.NewAuthMiddleware(userService)
	catBreedMiddleware := middleware.NewCatBreedMiddleware(catBreedService)

	// Настройка маршрутов
	http.HandleFunc("/api/register", userHandler.Register)
	http.HandleFunc("/api/login", userHandler.Login)

	// Защищенные маршруты пользователей
	http.Handle("/api/users", authMiddleware.RequireAuth(http.HandlerFunc(userHandler.GetAllUsers)))
	http.Handle("/api/user",
		authMiddleware.RequireAuth(
			authMiddleware.RequireUserAccessOrAdmin(
				http.HandlerFunc(userHandler.GetUser),
			),
		),
	)
	http.Handle("/api/user/update",
		authMiddleware.RequireAuth(
			authMiddleware.RequireUserAccessOrAdmin(
				http.HandlerFunc(userHandler.UpdateUser),
			),
		),
	)
	http.Handle("/api/user/delete",
		authMiddleware.RequireAuth(
			authMiddleware.RequireUserAccessOrAdmin(
				http.HandlerFunc(userHandler.DeleteUser),
			),
		),
	)

	// Маршруты для пород кошек
	// Публичные маршруты (не требуют аутентификации)
	http.HandleFunc("/api/cat-breeds", catBreedHandler.GetAllCatBreeds)
	http.HandleFunc("/api/cat-breed", catBreedHandler.GetCatBreed)

	// Защищенные маршруты (требуют аутентификации)
	http.Handle("/api/cat-breeds/my", authMiddleware.RequireAuth(http.HandlerFunc(catBreedHandler.GetUserCatBreeds)))
	http.Handle("/api/cat-breed/create", authMiddleware.RequireAuth(http.HandlerFunc(catBreedHandler.Create)))
	http.Handle("/api/cat-breed/update",
		authMiddleware.RequireAuth(
			catBreedMiddleware.RequireCatBreedOwnerOrAdmin(
				http.HandlerFunc(catBreedHandler.UpdateCatBreed),
			),
		),
	)
	http.Handle("/api/cat-breed/delete",
		authMiddleware.RequireAuth(
			catBreedMiddleware.RequireCatBreedOwnerOrAdmin(
				http.HandlerFunc(catBreedHandler.DeleteCatBreed),
			),
		),
	)

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
