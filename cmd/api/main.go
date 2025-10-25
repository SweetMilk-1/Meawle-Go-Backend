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

	// Инициализация сервисов
	userService := services.NewUserService(userRepo, jwtSecret)

	// Инициализация хэндлеров
	userHandler := handlers.NewUserHandler(userService)

	// Инициализация middleware
	authMiddleware := middleware.NewAuthMiddleware(userService)

	// Настройка маршрутов
	http.HandleFunc("/api/register", userHandler.Register)
	http.HandleFunc("/api/login", userHandler.Login)

	// Защищенные маршруты
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

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
